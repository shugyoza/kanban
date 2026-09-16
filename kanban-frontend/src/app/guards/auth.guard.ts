import { inject } from '@angular/core';
import { CanActivateFn, Router } from '@angular/router';
import { AuthService } from '../services/auth.service';
import { toObservable } from '@angular/core/rxjs-interop';
import { filter, map, take } from 'rxjs';

export const authGuard: CanActivateFn = () => {
  const authService = inject(AuthService);
  const router = inject(Router);

  return toObservable(authService.authStatus).pipe(
    filter(status => status !== 'loading'),
    take(1),
    map(() => {
      const user = authService.currentUser();
      if (!user) {
        router.navigate(['/', 'login']);
      }

      return !!user;
    })
  )
};
