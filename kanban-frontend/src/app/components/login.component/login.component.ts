import { Component, inject, signal } from '@angular/core';
import { AuthService } from '../../services/auth.service';
import { LoginCredentials } from '../../models/auth.model';
import { form, minLength, required, FormField } from '@angular/forms/signals';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';
import { finalize } from 'rxjs';

@Component({
  imports: [FormField],
  selector: 'app-login.component',
  styleUrl: './login.component.css',
  templateUrl: './login.component.html',
})
export class LoginComponent {
  private readonly authService = inject(AuthService);

  private readonly credentialsModel = signal<Required<LoginCredentials>>({
    username: '',
    password: ''
  });

  protected readonly loginForm = form(this.credentialsModel, schemaPath => {
    required(schemaPath.username, { message: 'Username is required' });
    required(schemaPath.password, { message: 'Password is required' });
    minLength(schemaPath.password, 8, { message: 'Password must be at least 8 characters long' })
  });

  protected readonly loading = signal<boolean>(false);

  protected handleSubmit($event: Event): void {
    $event.preventDefault();

    if (this.loginForm().invalid()) {
      console.error('invalid login form');

      return;
    }

    this.loading.set(true);
    this.authService.login(this.credentialsModel()).pipe(
      takeUntilDestroyed(),
      finalize(() => {
        this.loading.set(false)
      })
    ).subscribe();
  }
}
