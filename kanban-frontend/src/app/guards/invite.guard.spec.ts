import { TestBed } from '@angular/core/testing';
import { CanActivateFn } from '@angular/router';
import { inviteGuard } from './invite.guard';

describe('inviteGuard', () => {
  const executeGuard: CanActivateFn = (...guardParameters) =>
    TestBed.runInInjectionContext(() => inviteGuard(...guardParameters));

  beforeEach(() => {
    TestBed.configureTestingModule({});
  });

  it('should be created', () => {
    expect(executeGuard).toBeTruthy();
  });
});
