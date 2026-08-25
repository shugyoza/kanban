import { Component, signal } from '@angular/core';
import { RouterOutlet } from '@angular/router';

import { BoardComponent } from './board.component/board.component';

@Component({
  imports: [RouterOutlet, BoardComponent],
  selector: 'app-root',
  styleUrl: './app.css',
  templateUrl: './app.html',
})
export class App {
  protected readonly title = signal('kanban-frontend');
}
