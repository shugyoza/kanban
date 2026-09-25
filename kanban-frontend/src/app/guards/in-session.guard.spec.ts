import { TestBed } from '@angular/core/testing';
import { CanActivateFn } from '@angular/router';
import { inSessionGuard } from './in-session.guard';

describe('inSessionGuard', () => {
  const executeGuard: CanActivateFn = (...guardParameters) =>
    TestBed.runInInjectionContext(() => inSessionGuard(...guardParameters));

  beforeEach(() => {
    TestBed.configureTestingModule({});
  });

  it('should be created', () => {
    expect(executeGuard).toBeTruthy();
  });
});
