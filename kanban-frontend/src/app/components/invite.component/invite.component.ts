import { Component, DestroyRef, inject, signal } from '@angular/core';
import { InvitationService } from '../../services/invitation.service';
import { form, FormField, required } from '@angular/forms/signals';
import { AuthService } from '../../services/auth.service';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';

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
  private readonly emailModel = signal<Required<string>>('');

  protected readonly emailForm = form(this.emailModel, schemaPath => {
    required(schemaPath, { message: 'Email is required' })
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

    if (this.emailForm().invalid()) {
      console.error('invalid email form');

      return;
    }

    this.inviteService.getInvitationToken({
      userId,
      email: this.emailModel()
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
