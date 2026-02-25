'use client';

import { FileItem } from '@/types';

export const mockFiles: FileItem[] = [
  {
    id: '1',
    name: 'Documents',
    type: 'folder',
    modifiedAt: new Date('2024-01-15'),
    owner: 'me',
  },
  {
    id: '2',
    name: 'Photos',
    type: 'folder',
    modifiedAt: new Date('2024-02-01'),
    owner: 'me',
  },
  {
    id: '3',
    name: 'Projects',
    type: 'folder',
    modifiedAt: new Date('2024-01-20'),
    owner: 'me',
    shared: true,
  },
  {
    id: '4',
    name: 'Resume.pdf',
    type: 'file',
    mimeType: 'application/pdf',
    size: 245760,
    modifiedAt: new Date('2024-01-10'),
    owner: 'me',
    starred: true,
  },
  {
    id: '5',
    name: 'Budget 2024.xlsx',
    type: 'file',
    mimeType: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet',
    size: 35840,
    modifiedAt: new Date('2024-01-25'),
    owner: 'me',
  },
  {
    id: '6',
    name: 'Meeting Notes.docx',
    type: 'file',
    mimeType: 'application/vnd.openxmlformats-officedocument.wordprocessingml.document',
    size: 15360,
    modifiedAt: new Date('2024-01-28'),
    owner: 'me',
    shared: true,
  },
  {
    id: '7',
    name: 'Presentation.pptx',
    type: 'file',
    mimeType: 'application/vnd.openxmlformats-officedocument.presentationml.presentation',
    size: 1048576,
    modifiedAt: new Date('2024-02-05'),
    owner: 'me',
    starred: true,
  },
  {
    id: '8',
    name: 'vacation_photo.jpg',
    type: 'file',
    mimeType: 'image/jpeg',
    size: 2097152,
    modifiedAt: new Date('2024-01-05'),
    owner: 'me',
  },
  {
    id: '9',
    name: 'Notes.txt',
    type: 'file',
    mimeType: 'text/plain',
    size: 2048,
    modifiedAt: new Date('2024-02-10'),
    owner: 'me',
  },
  {
    id: '10',
    name: 'Project Plan.pdf',
    type: 'file',
    mimeType: 'application/pdf',
    size: 524288,
    modifiedAt: new Date('2024-02-08'),
    owner: 'john@example.com',
    shared: true,
  },
];

export const getFileIcon = (item: FileItem): string => {
  if (item.type === 'folder') return 'folder';
  
  const mimeType = item.mimeType || '';
  
  if (mimeType.includes('pdf')) return 'file-text';
  if (mimeType.includes('word') || mimeType.includes('document')) return 'file-text';
  if (mimeType.includes('sheet') || mimeType.includes('excel') || mimeType.includes('spreadsheet')) return 'file-spreadsheet';
  if (mimeType.includes('presentation') || mimeType.includes('powerpoint')) return 'presentation';
  if (mimeType.includes('image')) return 'image';
  if (mimeType.includes('video')) return 'video';
  if (mimeType.includes('audio')) return 'music';
  if (mimeType.includes('zip') || mimeType.includes('rar') || mimeType.includes('archive')) return 'archive';
  if (mimeType.includes('text')) return 'file-text';
  
  return 'file';
};

export const formatFileSize = (bytes?: number): string => {
  if (!bytes) return '—';
  
  const units = ['B', 'KB', 'MB', 'GB', 'TB'];
  let unitIndex = 0;
  let size = bytes;
  
  while (size >= 1024 && unitIndex < units.length - 1) {
    size /= 1024;
    unitIndex++;
  }
  
  return `${size.toFixed(1)} ${units[unitIndex]}`;
};

export const formatDate = (date: Date): string => {
  const now = new Date();
  const diff = now.getTime() - date.getTime();
  const days = Math.floor(diff / (1000 * 60 * 60 * 24));
  
  if (days === 0) return 'Today';
  if (days === 1) return 'Yesterday';
  if (days < 7) return `${days} days ago`;
  
  return date.toLocaleDateString('en-US', {
    month: 'short',
    day: 'numeric',
    year: date.getFullYear() !== now.getFullYear() ? 'numeric' : undefined,
  });
};
