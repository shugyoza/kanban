import { inject } from '@angular/core';
import { CanActivateFn, Router } from '@angular/router';
import { InvitationService } from '../services/invitation.service';
import { catchError, forkJoin, map, of } from 'rxjs';
import { HttpStatusCode } from '@angular/common/http';
import { AuthService } from '../services/auth.service';

export const inviteGuard: CanActivateFn = (route) => {
  const router = inject(Router);
  const inviteService = inject(InvitationService);
  const authService = inject(AuthService)

  const token = (route.queryParamMap.get('token') ?? '').trim();
  const email = (route.queryParamMap.get('email') ?? '').trim();

  if (!token || !email) {
    router.navigate(['/', 'login']);

    return of(false)
  }

  const isValidToken$ = inviteService.validateInvitationToken(token).pipe(
    map(response => response.status === HttpStatusCode.Ok),
    catchError(err => {
      console.error(err);

      return of(false)
    }),

  )
  const isRegisteredEmail$ = authService.validateEmailRegistered(email).pipe(
    map(response => response.status === HttpStatusCode.Ok && response.body?.registered === true),
    catchError(err => {
      console.error(err);

      return of(false)
    })
  )

  return forkJoin([
    isValidToken$,
    isRegisteredEmail$
  ]).pipe(
    map(([
      isValidToken,
      isRegisteredEmail
    ]) => {
      const emailCanBeRegistered = !isRegisteredEmail;

      if (!isValidToken || isRegisteredEmail) {
        router.navigate(['/', 'login']);
      }

      return isValidToken && emailCanBeRegistered;
    }),
  );

};
