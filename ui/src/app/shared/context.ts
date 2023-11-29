import { AppContext as WorkAppContext, NavigationApi, NotificationsApi, PopupApi } from 'work-ui';
import { History } from 'history';
import * as React from 'react';

export type AppContext = WorkAppContext & { apis: { popup: PopupApi; notifications: NotificationsApi; navigation: NavigationApi; baseHref: string } };

export interface ContextApis {
    popup: PopupApi;
    notifications: NotificationsApi;
    navigation: NavigationApi;
    history: History;
}

export const { Provider, Consumer } = React.createContext<ContextApis>(null);
