import { TestBed } from '@angular/core/testing';
import { KanbanService } from './kanban.service';
import { HttpTestingController, provideHttpClientTesting } from '@angular/common/http/testing';
import mockBoard from '../../../public/mock/board.json';
import { provideHttpClient } from '@angular/common/http';

describe('KanbanService', () => {
  let service: KanbanService;
  let httpMock: HttpTestingController;

  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [
        KanbanService,
        provideHttpClient(),
        provideHttpClientTesting(), // Configures the mock backend controller layer
      ]
    });

    service = TestBed.inject(KanbanService);
    httpMock = TestBed.inject(HttpTestingController);
  });

  afterEach(() => {
    // Verifies that there is no outstanding, unhandled HTTP network calls hanging around
    httpMock.verify();
  })

  it('should be created', () => {
    expect(service).toBeTruthy();
  });

  it('should initialize with a default null board state', () => {
    expect(service.board()).toBeNull();
    expect(service.isLoaded()).toEqual(false);
  })

  it('should fetch board tree details via proxy URL and update the reactive signal state', () => {
    const targetBoardId = 'board-1';

    // 1. Fire async network invocation
    service.loadBoard(targetBoardId);

    // 2. Intercept and evaluate the outgoing proxy request
    const http = httpMock.expectOne(`/api/boards?id=${targetBoardId}`);
    expect(http.request.method).toBe('GET');

    // 3. Flush mock dataset back down into the stream channel
    http.flush(mockBoard);

    // 4. Assert that writable Signal absorbed and updated the layout data models perfectly
    expect(service.board()).toEqual(mockBoard);
    expect(service.isLoaded()).toEqual(true);
    expect(service.board()?.title).toBe('Test Kanban Board')
  })

  it('should gracefully handle network failure streams and keep state safe', () => {
    const targetBoardId = 'invalid-board';

    // 1. Fire async network invocation
    service.loadBoard(targetBoardId);

    const http = httpMock.expectOne(`/api/boards?id=${targetBoardId}`);

    // Simulate a database/infrastructure connection drop error code from the server
    http.flush('Server Error', {
      status: 500,
      statusText: 'Internal Server Error'
    })

    expect(service.board()).toBeNull();
    expect(service.isLoaded()).toEqual(false);
  })
});
