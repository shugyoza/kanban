import { HttpClient } from '@angular/common/http';
import { computed, inject, Service, signal } from '@angular/core';
import { BoardAggregate, TaskCreateDTO, TaskUpdateDTO, Task, TaskEdit } from '../models/kanban.model';
import { catchError, of } from 'rxjs';
import { form, maxLength, required } from '@angular/forms/signals';

const INITIAL_TASK: TaskEdit = {
    title: '',
    description: ''
}

/* @Service() a new decorator in place of the old @Injectable. 
* By default @Service automatically registers the class as a global root-level singleton,
* similar to @Injectable({ providedIn: 'root' })
*/
@Service()
export class KanbanService {
    private readonly http = inject(HttpClient);
    private readonly boardState = signal<BoardAggregate | null>(null);
    public readonly board = this.boardState.asReadonly();
    public readonly isLoaded = computed<boolean>(() => this.boardState() !== null);
    public readonly taskIdOnEdit = signal<string | null>(null);
    public readonly isEditTaskFormOpen = signal<boolean>(false);
    public readonly editTaskModel = signal<TaskEdit>({ ...INITIAL_TASK });
    public readonly editTaskForm = form(this.editTaskModel, schemaPath => {
        required(schemaPath.title, { message: 'Title is required' });
        required(schemaPath.description, { message: 'Description is required' });
        maxLength(schemaPath.title, 255, { message: 'Maximum 255 characters' })
    });


    public loadBoard(boardId: string): void {
        this.http.get<BoardAggregate>(`/api/boards?id=${boardId}`).pipe(
            catchError(error => {
                console.error('Data stream resolution failed: ', error);

                return of(null)
            })
        ).subscribe(data => {
            this.boardState.set(data)
        })
    }

    public moveTask(
        column: {
            from: string;
            to: string;
        },
        row: {
            from: number;
            to: number;
        }
    ): void {
        const currentBoard = this.boardState();

        if (!currentBoard) return;

        // A. CAPTURE SNAPSHOT: Store a copy of the exact starting state in case server fails
        const rollbackSnapshot = { ...currentBoard };

        // 1. Deep copy the columns array to safely maintain immutability principles
        const updatedColumns = currentBoard.columns.map(c => ({
            ...c,
            tasks: [...c.tasks]
        }))

        const sourceColumn = updatedColumns.find(c => c.id === column.from);
        const targetColumn = updatedColumns.find(c => c.id === column.to);

        if (!sourceColumn || !targetColumn) return;

        // 2. Extract the target task being moved
        const [movedTask] = sourceColumn.tasks.splice(row.from, 1);
        if (!movedTask) return;

        // 3. Update the task's column reference identifier
        movedTask.columnId = column.to;

        // Inject the task into its new index position slot
        targetColumn.tasks.splice(row.to, 0, movedTask)

        // 5. Recalculate the sequential 'position' index property for both altered columns
        sourceColumn.tasks.forEach((task, index) => (task.position = index));
        if (column.from !== column.to) {
            targetColumn.tasks.forEach((task, index) => (task.position = index))
        }

        // 6. Push the updated model tree structural package into the state Signal.
        this.boardState.set({
            ...currentBoard,
            columns: updatedColumns
        })

        const payload = {
            taskId: movedTask.id,
            targetColumnId: column.to,
            targetPosition: row.to
        }

        this.http.put('/api/tasks/move', payload).pipe(
            catchError(err => {
                console.error(err);
                alert('Could not save card position. Checking database link...')

                // revert the state signal back
                this.boardState.set(rollbackSnapshot);

                return of(null);
            })
        ).subscribe();
    }

    public handleTaskEvent(columnId: string, $event: 'submit' | 'cancel'): void {
        if ($event === 'cancel') {
            this.isEditTaskFormOpen.set(false);
            this.taskIdOnEdit.set(null);
            this.editTaskForm().reset({ ...INITIAL_TASK });

            return;
        }

        const currentBoard = this.boardState();
        const title = this.editTaskModel().title;
        const description = this.editTaskModel().description;

        if (!currentBoard) return;

        // Capture snapshot to rollback current state if API call failed
        const rollbackSnapshot = { ...currentBoard };

        const taskId = this.taskIdOnEdit()
        const isTaskEdit = !!this.taskIdOnEdit();

        if (isTaskEdit) {

            // Deep copy columns and push existing cards down to clear index 0
            const updatedColumns = currentBoard.columns.map(col => {
                if (col.id !== columnId) {
                    return { ...col, tasks: [...col.tasks] };
                }

                const updatedTasks = col.tasks.map(t => {
                    if (t.id === taskId) {
                        return {
                            ...t,
                            title,
                            description
                        }
                    }

                    return {
                        ...t,
                    }
                })

                return {
                    ...col,
                    tasks: updatedTasks
                }
            })

            // Optimistic update
            this.boardState.set({
                ...currentBoard,
                columns: updatedColumns
            });

            const payload: TaskUpdateDTO = {
                taskId: taskId!,
                title,
                description
            }

            this.http.patch<Task>('/api/tasks', payload).pipe(
                catchError(error => {
                    console.error(error);
                    alert('Failed to save your changes. Reverting changes...');
                    this.boardState.set(rollbackSnapshot);

                    return of(null)
                })
            ).subscribe({
                next: () => {

                    // reset the create task form
                    this.editTaskForm().reset({ ...INITIAL_TASK });
                    this.isEditTaskFormOpen.set(false);
                    this.taskIdOnEdit.set(null);
                }
            })

            return;
        }

        // Create client-side id to trace the visual node
        const tempId = `temp-${Date.now()}`;

        // Deep copy columns and push existing cards down to clear index 0
        const updatedColumns = currentBoard.columns.map(col => {
            if (col.id !== columnId) {
                return { ...col, tasks: [...col.tasks] };
            }

            const shiftedTasks = col.tasks.map(t => {
                return {
                    ...t,
                    position: t.position + 1
                }
            })

            const newTask: Task = {
                id: tempId,
                columnId,
                title,
                description,
                position: 0
            }

            return {
                ...col,
                tasks: [newTask, ...shiftedTasks]
            }
        })

        // Optimistic update
        this.boardState.set({
            ...currentBoard,
            columns: updatedColumns
        })

        // Trigger Http call
        const payload: TaskCreateDTO = {
            columnId,
            title,
            description
        }

        this.http.post<Task>('/api/tasks', payload).pipe(
            catchError(error => {
                console.error(error);
                alert('Failed to save your new task. Reverting changes...');
                this.boardState.set(rollbackSnapshot);

                return of(null)
            })
        ).subscribe({
            next: serverSavedTask => {
                if (!serverSavedTask) return;

                // ID alignment: Swap out temporary task id with the final database ID returned from BE
                const finalizedBoard = this.boardState();
                if (!finalizedBoard) return;

                const alignedColumns = finalizedBoard.columns.map(col => {
                    if (col.id !== columnId) return { ...col, tasks: [...col.tasks] };

                    return {
                        ...col,
                        tasks: col.tasks.map(t => {
                            if (t.id === tempId) {

                                return {
                                    ...t,
                                    id: serverSavedTask.id
                                }
                            }

                            return t;
                        })
                    }
                })

                // finalized board state in UI
                this.boardState.set({
                    ...finalizedBoard, columns: alignedColumns
                })

                // reset the create task form
                this.editTaskForm().reset({ ...INITIAL_TASK });
                this.isEditTaskFormOpen.set(false);
            }
        })
    }

    public deleteTask(columnId: string, taskId: string, taskPosition: number): void {
        const currentBoard = this.boardState();
        if (!currentBoard) return;

        const rollbackSnapshot = { ...currentBoard };

        const updatedColumns = currentBoard.columns.map(col => {
            if (col.id !== columnId) {
                return { ...col, tasks: [...col.tasks] }
            }

            const filteredTasks = col.tasks.filter(t => t.id !== taskId);
            const reindexedTasks = filteredTasks.map((t, position) => {

                return {
                    ...t,
                    position
                }
            })

            return {
                ...col,
                tasks: reindexedTasks
            }
        })

        this.boardState.set({
            ...currentBoard,
            columns: updatedColumns
        })

        const options = {
            body: {
                columnId,
                taskId,
                taskPosition
            }
        }

        this.http.delete('/api/tasks', options).pipe(
            catchError(error => {
                console.error(error);
                this.boardState.set(rollbackSnapshot)
                return of(null)
            })
        ).subscribe()
    }

    public startEditingTask(task: Task): void {
        this.taskIdOnEdit.set(task.id);
        this.editTaskForm().reset({
            description: task.description,
            title: task.title
        })
    }

    public archiveTask(columnId: string, taskId: string, taskPosition: number): void {
        const currentBoard = this.boardState();
        if (!currentBoard) return;

        const rollbackSnapshot = { ...currentBoard };

        const updatedColumns = currentBoard.columns.map(col => {
            if (col.id !== columnId) {
                return { ...col, tasks: [...col.tasks] }
            }

            const filteredTasks = col.tasks.filter(t => t.id !== taskId);
            const reindexedTasks = filteredTasks.map((t, position) => {

                return {
                    ...t,
                    position
                }
            })

            return {
                ...col,
                tasks: reindexedTasks
            }
        })

        this.boardState.set({
            ...currentBoard,
            columns: updatedColumns
        })

        const payload = {
            columnId,
            taskId,
            taskPosition
        }

        this.http.patch('/api/tasks/archive', payload).pipe(
            catchError(error => {
                console.error(error);
                this.boardState.set(rollbackSnapshot)
                
                return of(null)
            })
        ).subscribe()
    }
}
