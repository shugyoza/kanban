import { HttpClient } from '@angular/common/http';
import { computed, inject, Service, signal } from '@angular/core';
import { Router } from '@angular/router';
import { LoginCredentials, User } from '../models/auth.model';
import { catchError, Observable, of, switchMap, tap, throwError } from 'rxjs';

@Service()
export class AuthService {
    private readonly http = inject(HttpClient);
    private readonly router = inject(Router);

    private currentUser$ = of<User | null>(null);
    private readonly currentUserState = signal<User | null>(null);
    public readonly currentUser = this.currentUserState.asReadonly();
    public readonly isAuthenticated = computed<boolean>(() => !!this.currentUser());

    public login(credentials: LoginCredentials): Observable<User | null> {
        return this.http.post<{ status: string }>(
            '/api/auth/login',
            credentials
        ).pipe(
            switchMap(response => {
                if (!!response && response.status === 'authenticated') {
                    this.currentUser$ = this.getCurrentUser();

                    return this.currentUser$;
                }

                this.currentUser$ = of(null);

                return this.currentUser$;
            }),
        )
    }

    public logout(): void {
        this.currentUserState.set(null);
        this.router.navigate(['/', 'login'])
    }

    public getCurrentUser(): Observable<User | null> {
        return this.http.get<User | null>(
            '/api/auth/me'
        ).pipe(
            tap(user => {
                this.currentUserState.set(user)
            }),
            catchError(error => {
                console.error(error);
                this.currentUserState.set(null);

                return throwError(() => error)
            })
        )
    }
}
