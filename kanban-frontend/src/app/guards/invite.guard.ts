import { inject } from '@angular/core';
import { CanActivateFn, Router } from '@angular/router';
import { InvitationService } from '../services/invitation.service';
import { catchError, map, of } from 'rxjs';
import { HttpStatusCode } from '@angular/common/http';

export const inviteGuard: CanActivateFn = (route) => {
  const router = inject(Router);
  const inviteService = inject(InvitationService);

  const token = route.queryParamMap.get('token');

  if (token && token.trim().length > 0) {
    return inviteService.validateInvitationToken(token).pipe(
      map(response => {
        const isValid = response.status === HttpStatusCode.Ok;
        if (!isValid) {
          router.navigate(['/', 'login']);
        }

        return isValid;
      }),
      catchError(error => {
        console.error(error);

        router.navigate(['/', 'login']);

        return of(false)
      })
    );
  }

  router.navigate(['/', 'login']);
  return false;
};
