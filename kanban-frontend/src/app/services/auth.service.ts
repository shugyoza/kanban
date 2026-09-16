import { HttpClient } from '@angular/common/http';
import { inject, Service, signal } from '@angular/core';
import { Router } from '@angular/router';
import { Credentials, User } from '../models/auth.model';
import { catchError, finalize, Observable, of, tap } from 'rxjs';
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

    public readonly currentUser = this.userResource.value.asReadonly();

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

    public logout(): void {
        this.http.delete<void>(
            '/api/auth/logout'
        ).pipe(
            finalize(() => {
                this.authVersion.set(0);
                this.router.navigate(['/', 'login'])
            })
        ).subscribe();
    }

    public getCurrentUser(): Observable<User | null> {
        return this.http.get<User | null>(
            '/api/auth/me'
        ).pipe(
            catchError(error => {
                console.error(error);

                return of(null)
            })
        )
    }
}
