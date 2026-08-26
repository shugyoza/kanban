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
        fromColumnId: string,
        toColumnId: string,
        fromIndex: number,
        toIndex: number
    ): void {
        const currentBoard = this.boardState();

        if (!currentBoard) return;

        // 1. Deep copy the columns array to safely maintain immutability principles
        const updatedColumns = currentBoard.columns.map(column => ({
            ...column,
            tasks: [...column.tasks]
        }))

        const sourceColumn = updatedColumns.find(column => column.id === fromColumnId);
        const targetColumn = updatedColumns.find(column => column);

        if (!sourceColumn || !targetColumn) return;

        // 2. Extract the target task being moved
        const [movedTask] = sourceColumn.tasks.splice(fromIndex, 1);
        if (!movedTask) return;

        // 3. Update the task's column reference identifier
        movedTask.columnId = toColumnId;

        // Inject the task into its new index position slot
        targetColumn.tasks.splice(toIndex, 0, movedTask)

        // 5. Recalculate the sequential 'position' index property for both altered columns
        sourceColumn.tasks.forEach((task, index) => (task.position = index));
        if (fromColumnId !== toColumnId) {
            targetColumn.tasks.forEach((task, index) => (task.position = index))
        }

        // 6. Push the updated model tree structural package into the state Signal.
        this.boardState.set({
            ...currentBoard,
            columns: updatedColumns
        })

        // 7. TODO: Trigger non blocking HTTP PATCH/PUT request to update the Go backend db persistence layer
        console.log(`Backend Synchronization Primed: Task ${movedTask.id} shifted to column ${toColumnId} at position ${toIndex}`)
    }
}
