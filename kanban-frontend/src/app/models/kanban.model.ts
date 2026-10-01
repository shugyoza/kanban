// should match with types in kanban-backend/internal/domain/kanban.go

export interface Task {
    id: string;
    columnId: string;
    title: string;
    description: string;
    position: number;
    isArchived?: boolean;
}

export interface Board {
    id: string;
    title: string;
}

export interface ColumnAggregate {
    id: string;
    title: string;
    position: number;
    tasks: Task[]
}

export interface BoardAggregate extends Board {
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