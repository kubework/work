import {AppContext} from 'work-ui/src/index';
import * as PropTypes from 'prop-types';
import * as React from 'react';
import {RouteComponentProps} from 'react-router';

export class BasePage<P extends RouteComponentProps<any>, S> extends React.Component<P, S> {
    public static contextTypes = {
        router: PropTypes.object,
        apis: PropTypes.object
    };

    public queryParam(name: string) {
        return this.params.get(name);
    }

    private get params() {
        return new URLSearchParams(this.appContext.router.route.location.search);
    }

    public queryParams(name: string) {
        return this.params.getAll(name);
    }

    // this allows us to set-multiple parameters at once
    public setQueryParams(newParams: any) {
        const params = this.params;
        Object.keys(newParams).forEach(name => {
            const value = newParams[name];
            if (value !== null) {
                params.set(name, value);
            } else {
                params.delete(name);
            }
        });
        this.pushParams(params);
    }

    public clearQueryParams() {
        this.appContext.router.history.push(this.props.match.url);
        this.refreshComponent();
    }

    // this allows us to set-multiple parameters at once
    public appendQueryParams(newParams: {name: string; value: string}[]) {
        const params = this.params;
        newParams.forEach(param => params.delete(param.name));
        newParams.forEach(param => params.append(param.name, param.value));
        this.pushParams(params);
    }

    private pushParams(params: URLSearchParams) {
        this.appContext.router.history.push(`${this.props.match.url}?${params.toString()}`);
        this.refreshComponent();
    }

    private refreshComponent() {
        setTimeout(() => {
            if (this.componentWillUnmount) {
                this.componentWillUnmount();
            }
            if (this.componentDidMount) {
                this.componentDidMount();
            }
        }, 300);
    }

    protected get appContext() {
        return this.context as AppContext;
    }
}
