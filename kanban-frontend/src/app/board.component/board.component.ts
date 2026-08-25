import { Component, inject, OnInit } from '@angular/core';
import { KanbanService } from '../services/kanban.service';

@Component({
  imports: [],
  selector: 'app-board',
  styleUrl: './board.component.css',
  templateUrl: './board.component.html',
})
export class BoardComponent implements OnInit {
  protected kanbanService = inject(KanbanService)

  ngOnInit(): void {
    this.kanbanService.loadBoard('board-1')
  }
}
