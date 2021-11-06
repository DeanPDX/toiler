export interface Task {
    id: number;
    userID: number;
    title: string;
    createdAt: Date;
    completedAt: Date;
}