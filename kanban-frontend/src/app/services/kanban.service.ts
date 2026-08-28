import { HttpClient } from '@angular/common/http';
import { computed, inject, Service, signal } from '@angular/core';
import { BoardAggregate, Task } from '../models/kanban.model';
import { catchError, of } from 'rxjs';

/* @Service() a new decorator in place of the old @Injectable. 
* By default @Service automatically registers the class as a global root-level singleton,
* similar to @Injectable({ providedIn: 'root' })
*/
@Service()
export class KanbanService {
    private http = inject(HttpClient);
    private boardState = signal<BoardAggregate | null>(null);
    public board = this.boardState.asReadonly();
    public isLoaded = computed<boolean>(() => this.boardState() !== null);
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

    public createTask(columnId: string, title: string, description: string):void {
        const currentBoard = this.boardState();
        if (!currentBoard) return;

        // Capture snapshot to rollback current state if API call failed
        const rollbackSnapshot = { ...currentBoard };
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
        const payload = {
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
                    if (col.id !== columnId) return col;

                    return {
                        ...col,
                        tasks: col.tasks.map(t => {
                            if (t.id !== tempId) {

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
            }
        })
    }
}
