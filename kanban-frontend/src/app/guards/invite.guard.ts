import { inject } from '@angular/core';
import { CanActivateFn, Router } from '@angular/router';

export const inviteGuard: CanActivateFn = (route) => {
  const router = inject(Router);
  const token = route.queryParamMap.get('token');

  if (token && token.trim().length > 0) {
    return true;
  }

  router.navigate(['/', 'login']);
  return false;
};
