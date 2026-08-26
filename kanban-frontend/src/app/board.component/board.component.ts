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

  // Method that handles pointer drag drop operations natively
  protected onCardDropped(event: CdkDragDrop<Task[]>): void {
    // If card was picked up and dropped back into its exact starting container slot index, skip processing
    if (event.previousContainer === event.container && event.previousIndex === event.currentIndex) {

      return;
    }

    console.log({
      fromIndex: event.previousIndex,
      toIndex: event.currentIndex,
      fromContainerData: event.previousContainer.data,
      toContainerData: event.container.data
    })
  }
}
