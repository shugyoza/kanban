import { HttpClient } from '@angular/common/http';
import { computed, inject, Service, signal } from '@angular/core';
import { BoardAggregate } from '../models/kanban.model';
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

        // TODO: clean up
        console.log({
            sourceColumn,
            targetColumn
        })

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

        // 7. TODO: Trigger non blocking HTTP PATCH/PUT request to update the Go backend db persistence layer
        console.log(`Backend Synchronization Primed: Task ${movedTask.id} shifted to column ${column.to} at position ${row.to}`)

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
}
