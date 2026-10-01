import { Component, computed, inject, signal } from '@angular/core';
import { KanbanService } from '../../services/kanban.service';
import { FormField } from '@angular/forms/signals';
import { finalize, take } from 'rxjs';

@Component({
  imports: [FormField],
  selector: 'app-add-board',
  styleUrl: './add-board.component.css',
  templateUrl: './add-board.component.html',
})
export class AddBoardComponent {
  protected readonly kanbanService = inject(KanbanService);
  protected readonly newBoardTitleForm = computed(() => this.kanbanService.newBoardTitleForm());
  protected readonly loading = signal<boolean>(false);
  protected submitNewBoardTitle($event: Event): void {
    $event.preventDefault();

    if (this.newBoardTitleForm().invalid()) {
      alert('invalid form');

      return;
    }

    this.loading.set(true);
    this.kanbanService.submitNewBoardTitle().pipe(
      take(1),
      finalize(() => this.loading.set(false))
    ).subscribe();
  }
}
