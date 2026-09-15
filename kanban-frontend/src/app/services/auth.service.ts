import { HttpClient } from '@angular/common/http';
import { computed, inject, Service, signal, WritableSignal } from '@angular/core';
import { Router } from '@angular/router';
import { Credentials, User } from '../models/auth.model';
import { catchError, Observable, of, switchMap, tap, throwError } from 'rxjs';
import { rxResource } from '@angular/core/rxjs-interop';

@Service()
export class AuthService {
    private readonly http = inject(HttpClient);
    private readonly router = inject(Router);

    private readonly authVersion = signal<number>(0);
    private readonly userResource = rxResource<User | null, { version: number }>({
        params: () => ({
            version: this.authVersion()
        }),
        stream: () => this.getCurrentUser(),
    });

    public readonly isAuthenticated = computed<boolean>(() => !!this.currentUser());

    public readonly currentUser: WritableSignal<User | null | undefined> = this.userResource.value;

    public register(credentials: Credentials): Observable<void> {
        return this.http.post<void>(
            'api/auth/register',
            credentials
        ).pipe(
            tap(() => {
                this.router.navigate(['/', 'login'])
            })
        )
    }

    public login(credentials: Credentials): Observable<void> {
        return this.http.post<void>(
            '/api/auth/login',
            credentials
        ).pipe(
            tap(() => {
                this.authVersion.update(v => v + 1)
            })
        )
    }

    public logout(): Observable<void> {
        return this.http.delete<void>(
            '/api/auth/logout'
        ).pipe(
            tap(() => {
                this.authVersion.update(v => v + 1);
                this.router.navigate(['/', 'login'])
            })
        )
    }

    public getCurrentUser(): Observable<User | null> {
        return this.http.get<User | null>(
            '/api/auth/me'
        )
    }
}
