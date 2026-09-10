import { Routes } from '@angular/router';
import { authGuard } from './guards/auth.guard';

export const routes: Routes = [
    {
        path: 'login',
        loadComponent: () => import('./components/login.component/login.component').then(m => m.LoginComponent)
    },
    {
        path: 'board',
        canActivate: [authGuard],
        loadComponent: () => import('./board.component/board.component').then(m => m.BoardComponent)
    },
    {
        path: '',
        redirectTo: 'board',
        pathMatch: 'full'
    },
    {
        path: '**',
        redirectTo: 'board'
    }
];
