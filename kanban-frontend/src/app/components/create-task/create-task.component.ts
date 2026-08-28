import { Component, inject, input, signal } from '@angular/core';
import { KanbanService } from '../../services/kanban.service';
import { FormField } from '@angular/forms/signals';

@Component({
  imports: [FormField],
  standalone: true,
  selector: 'app-create-task',
  styleUrl: './create-task.component.css',
  templateUrl: './create-task.component.html',
})
export class CreateTaskComponent {
  private kanbanService = inject(KanbanService);

  protected readonly createTaskForm = this.kanbanService.createTaskForm;
  protected readonly isCreateTaskFormOpen = this.kanbanService.isCreateTaskFormOpen;

  public columnId = input.required<string>();

  public submitTask(event: Event): void {
    event.preventDefault();
    this.kanbanService.createTask(this.columnId())
  }
}
