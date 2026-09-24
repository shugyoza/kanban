import { Routes } from '@angular/router';
import { authGuard } from './guards/auth.guard';
import { inviteGuard } from './guards/invite.guard';

export const routes: Routes = [
    {
        path: 'invite',
        canActivate: [authGuard],
        loadComponent: () => import('./components/invite.component/invite.component').then(m => m.InviteComponent)
    },
    {
        path: 'login',
        loadComponent: () => import('./components/login.component/login.component').then(m => m.LoginComponent)
    },
    {
        path: 'register',
        // canActivate: [inviteGuard], // TODO: re-instate once mailer implementations have been tested and verified
        loadComponent: () => import('./components/register.component/register.component').then(m => m.RegisterComponent)
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
