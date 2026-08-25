import { HttpClient } from '@angular/common/http';
import { computed, inject, Service, signal } from '@angular/core';
import { BoardAggregate } from '../models/kanban.model';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';
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
            takeUntilDestroyed(), // we don't really need to unsubscribe as this is a singleton provided in root, but this is for best practice
            catchError(error => {
                console.error('Data stream resolution failed: ', error);

                return of(null)
            })
        ).subscribe(data => {
            this.boardState.set(data)
        })
    }
}
