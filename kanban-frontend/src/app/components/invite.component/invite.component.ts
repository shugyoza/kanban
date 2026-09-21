import { Component, DestroyRef, inject, signal } from '@angular/core';
import { InvitationService } from '../../services/invitation.service';
import { debounce, email, form, FormField, required, validateHttp } from '@angular/forms/signals';
import { AuthService } from '../../services/auth.service';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';
import { HttpErrorResponse, HttpResponse, HttpStatusCode } from '@angular/common/http';

@Component({
  imports: [FormField],
  selector: 'app-invite.component',
  styleUrl: './invite.component.css',
  templateUrl: './invite.component.html',
})
export class InviteComponent {
  private readonly inviteService = inject(InvitationService);
  private readonly authService = inject(AuthService);
  private readonly destroyRef = inject(DestroyRef);

  private readonly user = this.authService.currentUser()
  private validatedEmails = new Set<string>();
  private readonly inviteModel = signal<Required<{ email: string }>>({
    email: '',
  });

  protected readonly inviteForm = form(this.inviteModel, schemaPath => {
    required(schemaPath.email, { message: 'Email is required' });
    email(schemaPath.email, { message: 'Invalid email'});

    debounce(schemaPath.email, 300)
    validateHttp(schemaPath.email, {
      request: ({ value }) => ({
        url: '/api/email/validate',
        method: 'POST',
        body: { email: value() },
      }),
      onSuccess: (response: HttpResponse<void>, { value }) => {
        if (response.ok) {
          // Cache successful validations
          this.validatedEmails.add(value());
        }

        return null;
      },
      onError: (err: unknown) => {
        const error = err as HttpErrorResponse;

        if (error.status === HttpStatusCode.BadRequest && error.error === 'Invalid email input') {

          return {
            kind: 'invalid',
            message: error.error
          }
        }
        
        return {
          kind: 'requestFailed',
          message: 'Unable to validate email input'
        }
      }
    })
  })
  protected readonly loading = signal<boolean>(false);
  protected readonly invitationURL = signal<string>('');

  protected submitInvitationTokenRequest($event: Event): void {
    $event.preventDefault();

    const userId = this.user?.id;
    if (!userId) {
      console.error('need to login')

      return;
    }

    if (this.inviteForm.email().invalid()) {
      console.error('invalid email form');

      return;
    }

    this.inviteService.getInvitationToken({
      userId,
      email: this.inviteModel().email
    }).pipe(
      takeUntilDestroyed(this.destroyRef),
    ).subscribe({
      next: response => {
        const { token } = response;
        const baseHref = window.location.origin;

        const invitationURL = `${baseHref}/register?token=${token}`;
        this.invitationURL.set(invitationURL);
      },
      error: error => {
        console.error(error);
        this.invitationURL.set('')
      }
    })
  }

  public async copyInvitationURL(): Promise<void> {
    const url = this.invitationURL();

    try {
      await navigator.clipboard.writeText(url);
      alert('Invitation URL copied!')
    } catch (err) {
      console.error('Could not copy invitation URL: ', err)
    }
  }
}
