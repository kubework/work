import { Page } from 'work-ui';
import * as React from 'react';

require('./help.scss');

export const Help = () => (
    <Page title='Help'>
        <div className='row'>
            <div className='columns large-4 medium-12'>
                <div className='help-box'>
                    <div className='help-box__ico help-box__ico--manual' />
                    <h3>Documentation</h3>
                    <a href='https://kubework.github.io/' target='_blank' className='help-box__link'>
                        {' '}
                        Work Project
                    </a>
                </div>
            </div>
            <div className='columns large-4 medium-12'>
                <div className='help-box'>
                    <div className='help-box__ico help-box__ico--email' />
                    <h3>Contact</h3>
                    <a className='help-box__link' target='_blank' href='https://groups.google.com/forum/#!forum/kubework'>
                        Work Community
                    </a>
                    <a className='help-box__link' target='_blank' href='https://kubework.slack.com'>
                        Slack Channel
                    </a>
                </div>
            </div>
            <div className='columns large-4 medium-12'>
                <div className='help-box'>
                    <div className='help-box__ico help-box__ico--download' />
                    <h3>Work CLI</h3>
                    <div className='row text-left help-box__download'>
                        <div className='columns small-4'>
                            <a href={`https://github.com/kubework/work/releases/download/${SYSTEM_INFO.version}/work-linux-amd64`}>
                                <i className='fab fa-linux' aria-hidden='true' /> Linux
                            </a>
                        </div>
                        <div className='columns small-4'>
                            <a href={`https://github.com/kubework/work/releases/download/${SYSTEM_INFO.version}/work-darwin-amd64`}>
                                <i className='fab fa-apple' aria-hidden='true' /> macOS
                            </a>
                            <br />
                        </div>
                        <div className='columns small-4'>
                            <a href={`https://github.com/kubework/work/releases/download/${SYSTEM_INFO.version}/work-windows-amd64`}>
                                <i className='fab fa-windows' aria-hidden='true' /> Windows
                            </a>
                            <br />
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </Page>
);
