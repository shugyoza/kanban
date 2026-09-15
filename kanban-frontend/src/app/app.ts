import { Component, computed, inject, signal } from '@angular/core';
import { RouterOutlet } from '@angular/router';
import { AuthService } from './services/auth.service';

@Component({
  imports: [RouterOutlet],
  selector: 'app-root',
  styleUrl: './app.css',
  templateUrl: './app.html',
})
export class App {
  private readonly authService = inject(AuthService);
  protected readonly title = signal('kanban-frontend');
  protected readonly user = computed(() => this.authService.currentUser());
}
