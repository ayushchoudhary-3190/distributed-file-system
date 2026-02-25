'use client';

import { FileItem } from '@/types';
import FileItemCard from './FileItem';
import { useState } from 'react';

interface FileGridProps {
  files: FileItem[];
  viewMode: 'grid' | 'list';
}

export default function FileGrid({ files, viewMode }: FileGridProps) {
  const [selectedFiles, setSelectedFiles] = useState<Set<string>>(new Set());

  const handleSelect = (id: string) => {
    setSelectedFiles(prev => {
      const newSet = new Set(prev);
      if (newSet.has(id)) {
        newSet.delete(id);
      } else {
        newSet.add(id);
      }
      return newSet;
    });
  };

  const handleStar = (id: string) => {
    console.log('Star:', id);
  };

  const handleClick = (id: string) => {
    console.log('Click:', id);
  };

  if (files.length === 0) {
    return (
      <div className="flex-1 flex items-center justify-center">
        <div className="text-center">
          <div className="w-24 h-24 mx-auto mb-4 bg-white/5 rounded-full flex items-center justify-center">
            <span className="text-4xl">📁</span>
          </div>
          <p className="text-gray-400 text-lg">No files in this view</p>
          <p className="text-gray-500 text-sm mt-1">Drop files here or click New to upload</p>
        </div>
      </div>
    );
  }

  if (viewMode === 'list') {
    return (
      <div className="flex-1 overflow-auto">
        <div className="flex items-center gap-4 px-4 py-2 border-b border-white/10 text-gray-400 text-sm sticky top-0 bg-[#1f1f1f]">
          <div className="w-5" />
          <div className="w-10" />
          <div className="flex-1">Name</div>
          <div className="w-32">Owner</div>
          <div className="w-32">Last modified</div>
          <div className="w-24 text-right">Size</div>
          <div className="w-8" />
        </div>
        {files.map((file) => (
          <FileItemCard
            key={file.id}
            item={file}
            viewMode={viewMode}
            onClick={() => handleClick(file.id)}
            onStar={() => handleStar(file.id)}
            isSelected={selectedFiles.has(file.id)}
            onSelect={() => handleSelect(file.id)}
          />
        ))}
      </div>
    );
  }

  return (
    <div className="flex-1 overflow-auto p-4">
      <div className="grid grid-cols-[repeat(auto-fill,minmax(180px,1fr))] gap-4">
        {files.map((file) => (
          <FileItemCard
            key={file.id}
            item={file}
            viewMode={viewMode}
            onClick={() => handleClick(file.id)}
            onStar={() => handleStar(file.id)}
            isSelected={selectedFiles.has(file.id)}
            onSelect={() => handleSelect(file.id)}
          />
        ))}
      </div>
    </div>
  );
}
