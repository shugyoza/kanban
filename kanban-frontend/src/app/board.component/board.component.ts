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
    this.kanbanService.loadBoard('board-kanban-1')
  }

  // Method that handles pointer drag drop operations natively
  protected onCardDropped(event: CdkDragDrop<Task[]>): void {
    const column = {
      from: event.previousContainer.id,
      to: event.container.id
    }
    const row = {
      from: event.previousIndex,
      to: event.currentIndex
    }
    const data = {
      from: event.previousContainer.data,
      to: event.container.data
    }
    // If card was picked up and dropped back into its exact starting container slot index, skip processing
    if (column.from === column.to && row.from === row.to) {

      return;
    }

    // Log movements for developing purpose. TODO: clean up later
    console.log(column, row, data)

    // Pass the indices down to state store orchestrator service layer
    this.kanbanService.moveTask(
      column,
      row
    )
  }
}
