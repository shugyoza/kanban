import { inject } from '@angular/core';
import { CanActivateFn, Router } from '@angular/router';
import { AuthService } from '../services/auth.service';
import { map, tap } from 'rxjs';

export const authGuard: CanActivateFn = () => {
  const authService = inject(AuthService);
  const router = inject(Router);

  return authService.getCurrentUser().pipe(
    map(user => !!user),
    tap(isLoggedIn => {
      if (!isLoggedIn) {
        router.navigate(['/', 'login']);
      }
    })
  )
};
