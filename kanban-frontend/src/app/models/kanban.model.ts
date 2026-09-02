// should match with types in kanban-backend/internal/domain/kanban.go

export interface Task {
    id: string;
    columnId: string;
    title: string;
    description: string;
    position: number;
}

export interface ColumnAggregate {
    id: string;
    title: string;
    position: number;
    tasks: Task[]
}

export interface BoardAggregate {
    id: string;
    title: string;
    columns: ColumnAggregate[];
}

export interface TaskEdit {
  title: string;
  description: string;
}

export interface TaskCreateDTO extends TaskEdit {
  columnId: string;
}

export interface TaskUpdateDTO extends TaskEdit {
  taskId: string;
}