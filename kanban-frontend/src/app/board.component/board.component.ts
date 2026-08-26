import { Component, inject, OnInit } from '@angular/core';
import { KanbanService } from '../services/kanban.service';
import { Task } from '../models/kanban.model';
import { CdkDragDrop, DragDropModule } from '@angular/cdk/drag-drop';

@Component({
  imports: [DragDropModule],
  selector: 'app-board',
  styleUrl: './board.component.css',
  templateUrl: './board.component.html',
})
export class BoardComponent implements OnInit {
  protected kanbanService = inject(KanbanService)

  ngOnInit(): void {
    this.kanbanService.loadBoard('board-1')
  }

  protected onCardDropped(event: CdkDragDrop<Task[]>): void {

  }
}
