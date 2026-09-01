import { Component, inject, output } from '@angular/core';
import { KanbanService } from '../../services/kanban.service';
import { TaskCreateDTO } from '../../models/kanban.model';
import { FormField } from '@angular/forms/signals';

@Component({
  imports: [FormField],
  standalone: true,
  selector: 'app-edit-task',
  styleUrl: './edit-task.component.css',
  templateUrl: './edit-task.component.html',
})
export class EditTaskComponent {
  private kanbanService = inject(KanbanService);
  protected readonly editTaskForm = this.kanbanService.editTaskForm;
  public readonly taskEvent = output<'submit' | 'cancel'>();
}
