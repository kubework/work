package sqldb

import (
	"context"

	log "github.com/sirupsen/logrus"
	"upper.io/db.v3/lib/sqlbuilder"
)

type Migrate interface {
	Exec(ctx context.Context) error
}

func NewMigrate(session sqlbuilder.Database, clusterName string, tableName string) Migrate {
	return migrate{session, clusterName, tableName}
}

type migrate struct {
	session     sqlbuilder.Database
	clusterName string
	tableName   string
}

type change interface {
	apply(session sqlbuilder.Database) error
}

func ternary(condition bool, left, right change) change {
	if condition {
		return left
	} else {
		return right
	}
}

func (m migrate) Exec(ctx context.Context) error {
	{
		// poor mans SQL migration
		_, err := m.session.Exec("create table if not exists schema_history(schema_version int not null)")
		if err != nil {
			return err
		}
		rs, err := m.session.Query("select schema_version from schema_history")
		if err != nil {
			return err
		}
		if !rs.Next() {
			_, err := m.session.Exec("insert into schema_history values(-1)")
			if err != nil {
				return err
			}
		}
		err = rs.Close()
		if err != nil {
			return err
		}
	}
	dbType := dbTypeFor(m.session)

	log.WithFields(log.Fields{"clusterName": m.clusterName, "dbType": dbType}).Info("Migrating database schema")

	// try and make changes idempotent, as it is possible for the change to apply, but the archive update to fail
	// and therefore try and apply again next try

	for changeSchemaVersion, change := range []change{
		ansiSQLChange(`create table if not exists ` + m.tableName + ` (
    id varchar(128) ,
    name varchar(256),
    phase varchar(25),
    namespace varchar(256),
    workflow text,
    startedat timestamp,
    finishedat timestamp,
    primary key (id, namespace)
)`),
		ansiSQLChange(`create unique index idx_name on ` + m.tableName + ` (name)`),
		ansiSQLChange(`create table if not exists work_workflow_history (
    id varchar(128) ,
    name varchar(256),
    phase varchar(25),
    namespace varchar(256),
    workflow text,
    startedat timestamp,
    finishedat timestamp,
    primary key (id, namespace)
)`),
		ansiSQLChange(`alter table work_workflow_history rename to work_archived_workflows`),
		ternary(dbType == MySQL,
			ansiSQLChange(`drop index idx_name on `+m.tableName),
			ansiSQLChange(`drop index idx_name`),
		),
		ansiSQLChange(`create unique index idx_name on ` + m.tableName + `(name, namespace)`),
		ternary(dbType == MySQL,
			ansiSQLChange(`alter table `+m.tableName+` drop primary key`),
			ansiSQLChange(`alter table `+m.tableName+` drop constraint `+m.tableName+`_pkey`),
		),
		ansiSQLChange(`alter table ` + m.tableName + ` add primary key(name,namespace)`),
		// huh - why does the pkey not have the same name as the table - history
		ternary(dbType == MySQL,
			ansiSQLChange(`alter table work_archived_workflows drop primary key`),
			ansiSQLChange(`alter table work_archived_workflows drop constraint work_workflow_history_pkey`),
		),
		ansiSQLChange(`alter table work_archived_workflows add primary key(id)`),
		// ***
		// THE CHANGES ABOVE THIS LINE MAY BE IN PER-PRODUCTION SYSTEMS - DO NOT CHANGE THEM
		// ***
		ansiSQLChange(`alter table work_archived_workflows rename column id to uid`),
		ternary(dbType == MySQL,
			ansiSQLChange(`alter table work_archived_workflows modify column uid varchar(128) not null`),
			ansiSQLChange(`alter table work_archived_workflows alter column uid set not null`),
		),
		ternary(dbType == MySQL,
			ansiSQLChange(`alter table work_archived_workflows modify column phase varchar(25) not null`),
			ansiSQLChange(`alter table work_archived_workflows alter column phase set not null`),
		),
		ternary(dbType == MySQL,
			ansiSQLChange(`alter table work_archived_workflows modify column namespace varchar(256) not null`),
			ansiSQLChange(`alter table work_archived_workflows alter column namespace set not null`),
		),
		ternary(dbType == MySQL,
			ansiSQLChange(`alter table work_archived_workflows modify column workflow text not null`),
			ansiSQLChange(`alter table work_archived_workflows alter column workflow set not null`),
		),
		ternary(dbType == MySQL,
			ansiSQLChange(`alter table work_archived_workflows modify column startedat timestamp not null`),
			ansiSQLChange(`alter table work_archived_workflows alter column startedat set not null`),
		),
		ternary(dbType == MySQL,
			ansiSQLChange(`alter table work_archived_workflows modify column finishedat timestamp not null`),
			ansiSQLChange(`alter table work_archived_workflows alter column finishedat set not null`),
		),
		ansiSQLChange(`alter table work_archived_workflows add clustername varchar(64)`), // DNS entry can only be max 63 bytes
		ansiSQLChange(`update work_archived_workflows set clustername = '` + m.clusterName + `' where clustername is null`),
		ternary(dbType == MySQL,
			ansiSQLChange(`alter table work_archived_workflows modify column clustername varchar(64) not null`),
			ansiSQLChange(`alter table work_archived_workflows alter column clustername set not null`),
		),
		ternary(dbType == MySQL,
			ansiSQLChange(`alter table work_archived_workflows drop primary key`),
			ansiSQLChange(`alter table work_archived_workflows drop constraint work_archived_workflows_pkey`),
		),
		ansiSQLChange(`alter table work_archived_workflows add primary key(clustername,uid)`),
		ansiSQLChange(`create index work_archived_workflows_i1 on work_archived_workflows (clustername,namespace)`),
		// work_archived_workflows now looks like:
		// clustername(not null) | uid(not null) | | name (null) | phase(not null) | namespace(not null) | workflow(not null) | startedat(not null)  | finishedat(not null)
		// remove unused columns
		ansiSQLChange(`alter table ` + m.tableName + ` drop column phase`),
		ansiSQLChange(`alter table ` + m.tableName + ` drop column startedat`),
		ansiSQLChange(`alter table ` + m.tableName + ` drop column finishedat`),
		ansiSQLChange(`alter table ` + m.tableName + ` rename column id to uid`),
		ternary(dbType == MySQL,
			ansiSQLChange(`alter table `+m.tableName+` modify column uid varchar(128) not null`),
			ansiSQLChange(`alter table `+m.tableName+` alter column uid set not null`),
		),
		ternary(dbType == MySQL,
			ansiSQLChange(`alter table `+m.tableName+` modify column namespace varchar(256) not null`),
			ansiSQLChange(`alter table `+m.tableName+` alter column namespace set not null`),
		),
		ansiSQLChange(`alter table ` + m.tableName + ` add column clustername varchar(64)`), // DNS cannot be longer than 64 bytes
		ansiSQLChange(`update ` + m.tableName + ` set clustername = '` + m.clusterName + `' where clustername is null`),
		ternary(dbType == MySQL,
			ansiSQLChange(`alter table `+m.tableName+` modify column clustername varchar(64) not null`),
			ansiSQLChange(`alter table `+m.tableName+` alter column clustername set not null`),
		),
		ansiSQLChange(`alter table ` + m.tableName + ` add column version varchar(64)`),
		ansiSQLChange(`alter table ` + m.tableName + ` add column nodes text`),
		backfillNodes{tableName: m.tableName},
		ternary(dbType == MySQL,
			ansiSQLChange(`alter table `+m.tableName+` modify column nodes text not null`),
			ansiSQLChange(`alter table `+m.tableName+` alter column nodes set not null`),
		),
		ansiSQLChange(`alter table ` + m.tableName + ` drop column workflow`),
		// add a timestamp column to indicate updated time
		ansiSQLChange(`alter table ` + m.tableName + ` add column updatedat timestamp not null default current_timestamp`),
		// remove the old primary key and add a new one
		ternary(dbType == MySQL,
			ansiSQLChange(`alter table `+m.tableName+` drop primary key`),
			ansiSQLChange(`alter table `+m.tableName+` drop constraint `+m.tableName+`_pkey`),
		),
		ternary(dbType == MySQL,
			ansiSQLChange(`drop index idx_name on `+m.tableName),
			ansiSQLChange(`drop index idx_name`),
		),
		ansiSQLChange(`alter table ` + m.tableName + ` drop column name`),
		ansiSQLChange(`alter table ` + m.tableName + ` add primary key(clustername,uid,version)`),
		ansiSQLChange(`create index ` + m.tableName + `_i1 on ` + m.tableName + ` (clustername,namespace)`),
		// work_workflows now looks like:
		//  clustername(not null) | uid(not null) | namespace(not null) | version(not null) | nodes(not null) | updatedat(not null)
		ternary(dbType == MySQL,
			ansiSQLChange(`alter table work_archived_workflows modify column workflow json not null`),
			ansiSQLChange(`alter table work_archived_workflows alter column workflow type json using workflow::json`),
		),
		ternary(dbType == MySQL,
			ansiSQLChange(`alter table work_archived_workflows modify column name varchar(256) not null`),
			ansiSQLChange(`alter table work_archived_workflows alter column name set not null`),
		),
		// clustername(not null) | uid(not null) | | name (not null) | phase(not null) | namespace(not null) | workflow(not null) | startedat(not null)  | finishedat(not null)
		ansiSQLChange(`create index ` + m.tableName + `_i2 on ` + m.tableName + ` (clustername,namespace,updatedat)`),
		// The work_archived_workflows_labels is really provided as a way to create queries on labels that are fast because they
		// use indexes. When displaying, it might be better to look at the `workflow` column.
		// We could have added a `labels` column to work_archived_workflows, but then we would have had to do free-text
		// queries on it which would be slow due to having to table scan.
		// The key has an optional prefix(253 chars) + '/' + name(63 chars)
		// Why is the key called "name" not "key"? Key is an SQL reserved word.
		ansiSQLChange(`create table if not exists work_archived_workflows_labels (
	clustername varchar(64) not null,
	uid varchar(128) not null,
    name varchar(317) not null,
    value varchar(63) not null,
    primary key (clustername, uid, name),
 	foreign key (clustername, uid) references work_archived_workflows(clustername, uid) on delete cascade
)`),
		// MySQL can only store 64k in a TEXT field, both MySQL and Posgres can store 1GB in JSON.
		ternary(dbType == MySQL,
			ansiSQLChange(`alter table `+m.tableName+` modify column nodes json not null`),
			ansiSQLChange(`alter table `+m.tableName+` alter column nodes type json using nodes::json`),
		),
	} {
		err := m.applyChange(ctx, changeSchemaVersion, change)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m migrate) applyChange(ctx context.Context, changeSchemaVersion int, c change) error {
	tx, err := m.session.NewTx(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()
	rs, err := tx.Exec("update schema_history set schema_version = ? where schema_version = ?", changeSchemaVersion, changeSchemaVersion-1)
	if err != nil {
		return err
	}
	rowsAffected, err := rs.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 1 {
		log.WithFields(log.Fields{"changeSchemaVersion": changeSchemaVersion, "change": c}).Info("applying database change")
		err := c.apply(m.session)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}
