import { Page } from 'work-ui';
import * as React from 'react';
import { uiUrl } from '../../shared/base';

require('./login.scss');

const getToken = () => {
    for (const cookie of document.cookie.split(';')) {
        if (cookie.startsWith('authorization=')) {
            return cookie.substring(14);
        }
    }
    return null;
};

const maybeLoggedIn = () => !!getToken();
const logout = () => {
    document.cookie = 'authorization=;';
    document.location.reload(true);
};
const login = (token: string) => {
    document.cookie = 'authorization=' + token + ';';
    document.location.href = uiUrl('');
};
export const Login = () => (
    <Page title='Login' toolbar={{ breadcrumbs: [{ title: 'Login' }] }}>
        <div className='work-container'>
            <p>
                <i className='fa fa-info-circle' /> You appear to be logged {maybeLoggedIn() ? 'in' : 'out'}. It may not be necessary to login to use Work, it depends on how it is
                configured.
            </p>
            <p>
                Get your token using <code>work auth token</code> and paste in this box.
            </p>
            <textarea id='token' cols={100} rows={20} defaultValue={getToken()} />
            <div>
                {maybeLoggedIn() && (
                    <button className='work-button work-button--base-o' onClick={() => logout()}>
                        <i className='fa fa-lock' /> Logout
                    </button>
                )}
                <button className='work-button work-button--base-o' onClick={() => login((document.getElementById('token') as HTMLInputElement).value)}>
                    <i className='fa fa-lock-open' /> Login
                </button>
            </div>
        </div>
    </Page>
);
