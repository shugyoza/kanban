import { Component, inject, signal } from '@angular/core';
import { AuthService } from '../../services/auth.service';
import { Credentials } from '../../models/auth.model';
import { ActivatedRoute, Router } from '@angular/router';
import { form, FormField, minLength, required } from '@angular/forms/signals';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';
import { finalize } from 'rxjs';

@Component({
  imports: [FormField],
  selector: 'app-register.component',
  styleUrl: './register.component.css',
  templateUrl: './register.component.html',
})
export class RegisterComponent {
  private readonly authService = inject(AuthService);
  private readonly router = inject(Router);
  private readonly activatedRoute = inject(ActivatedRoute);
  private readonly token = this.activatedRoute.snapshot.queryParamMap.get('token') ?? '';

  private readonly credentialsModel = signal<Required<Credentials>>({
    username: '',
    password: ''
  })

  protected readonly registerForm = form(this.credentialsModel, schemaPath => {
    required(schemaPath.username, { message: 'Username is required' });
    required(schemaPath.password, { message: 'Password is required' });
    minLength(schemaPath.password, 8, { message: 'Password must be at least 8 characters long' })
  });

  protected readonly loading = signal<boolean>(false);

  protected submitRegister($event: Event): void {
    $event.preventDefault();

    if (this.registerForm().invalid()) {
      console.error('invalid register form');

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
