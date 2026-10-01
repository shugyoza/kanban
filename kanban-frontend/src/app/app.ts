import { Component, computed, inject, signal } from '@angular/core';
import { RouterOutlet } from '@angular/router';
import { AuthService } from './services/auth.service';
import { User } from './models/auth.model';
import { AddBoardComponent } from './components/add-board.component/add-board.component';
import { BoardsComponent } from './components/boards.component/boards.component';

@Component({
  imports: [RouterOutlet, AddBoardComponent, BoardsComponent],
  selector: 'app-root',
  styleUrl: './app.css',
  templateUrl: './app.html',
})
export class App {
  private readonly authService = inject(AuthService);
  protected readonly title = signal('kanban-frontend');
  protected readonly user = computed<User | null | undefined>(() => this.authService.currentUser());

  public logout(): void {
    this.authService.logout();
  }
}
