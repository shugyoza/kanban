import { Component, DestroyRef, inject, signal } from '@angular/core';
import { AuthService } from '../../services/auth.service';
import { RegisterCredentials } from '../../models/auth.model';
import { ActivatedRoute } from '@angular/router';
import { disabled, form, FormField, minLength, required } from '@angular/forms/signals';
import { finalize } from 'rxjs';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';

@Component({
  imports: [FormField],
  selector: 'app-register.component',
  styleUrl: './register.component.css',
  templateUrl: './register.component.html',
})
export class RegisterComponent {
  private readonly authService = inject(AuthService);
  private readonly activatedRoute = inject(ActivatedRoute);
  private readonly destroyRef = inject(DestroyRef);
  private readonly token = (this.activatedRoute.snapshot.queryParamMap.get('token') ?? '').trim();
  private readonly email = (this.activatedRoute.snapshot.queryParamMap.get('email') ?? '').trim();

  private readonly credentialsModel = signal<Required<RegisterCredentials>>({
    username: '',
    password: '',
    email: this.email,
  })

  protected readonly registerForm = form(this.credentialsModel, schemaPath => {
    disabled(schemaPath.email)
    required(schemaPath.username, { message: 'Username is required' });
    required(schemaPath.password, { message: 'Password is required' });
    minLength(schemaPath.password, 8, { message: 'Password must be at least 8 characters long' });
  });

  protected readonly loading = signal<boolean>(false);

  protected submitRegister($event: Event): void {
    $event.preventDefault();

    if (this.registerForm().invalid()) {
      console.error('invalid register form');

      return;
    }

    this.loading.set(true);
    this.authService.register(this.credentialsModel()).pipe(
      takeUntilDestroyed(this.destroyRef),
      finalize(() => {
        this.loading.set(false)
      })
    ).subscribe();
  }
}
