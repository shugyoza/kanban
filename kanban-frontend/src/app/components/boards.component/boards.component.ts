import { Component, inject } from '@angular/core';
import { KanbanService } from '../../services/kanban.service';
import { AsyncPipe, JsonPipe } from '@angular/common';

@Component({
  imports: [AsyncPipe, JsonPipe],
  selector: 'app-boards',
  styleUrl: './boards.component.css',
  templateUrl: './boards.component.html',
})
export class BoardsComponent {
  private kanbanService = inject(KanbanService);

  protected readonly boards = this.kanbanService.boardsList;

  protected selectBoard(boardId: string): void {
    this.kanbanService.boardId.set(boardId);
  }
}
