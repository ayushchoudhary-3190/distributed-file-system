'use client';

import { 
  Folder, 
  FileText, 
  FileSpreadsheet, 
  Image, 
  Video, 
  Music, 
  Archive, 
  File,
  Presentation,
  Star,
  MoreVertical,
  Share2,
  Download,
  Pencil,
  Trash2,
  Copy
} from 'lucide-react';
import { useState } from 'react';
import { FileItem } from '@/types';
import { formatFileSize, formatDate } from '@/data/files';

interface FileItemCardProps {
  item: FileItem;
  viewMode: 'grid' | 'list';
  onClick: () => void;
  onStar: () => void;
  isSelected: boolean;
  onSelect: () => void;
}

const getIcon = (item: FileItem) => {
  if (item.type === 'folder') return <Folder className="w-10 h-10 text-yellow-400" />;
  
  const mimeType = item.mimeType || '';
  
  if (mimeType.includes('pdf') || mimeType.includes('word') || mimeType.includes('document') || mimeType.includes('text')) {
    return <FileText className="w-10 h-10 text-blue-400" />;
  }
  if (mimeType.includes('sheet') || mimeType.includes('excel') || mimeType.includes('spreadsheet')) {
    return <FileSpreadsheet className="w-10 h-10 text-green-400" />;
  }
  if (mimeType.includes('presentation') || mimeType.includes('powerpoint')) {
    return <Presentation className="w-10 h-10 text-orange-400" />;
  }
  if (mimeType.includes('image')) return <Image className="w-10 h-10 text-purple-400" />;
  if (mimeType.includes('video')) return <Video className="w-10 h-10 text-red-400" />;
  if (mimeType.includes('audio')) return <Music className="w-10 h-10 text-pink-400" />;
  if (mimeType.includes('zip') || mimeType.includes('rar') || mimeType.includes('archive')) {
    return <Archive className="w-10 h-10 text-yellow-600" />;
  }
  
  return <File className="w-10 h-10 text-gray-400" />;
};

const getIconBg = (item: FileItem): string => {
  if (item.type === 'folder') return 'bg-yellow-400/10';
  
  const mimeType = item.mimeType || '';
  
  if (mimeType.includes('pdf') || mimeType.includes('word') || mimeType.includes('document') || mimeType.includes('text')) {
    return 'bg-blue-400/10';
  }
  if (mimeType.includes('sheet') || mimeType.includes('excel') || mimeType.includes('spreadsheet')) {
    return 'bg-green-400/10';
  }
  if (mimeType.includes('presentation') || mimeType.includes('powerpoint')) {
    return 'bg-orange-400/10';
  }
  if (mimeType.includes('image')) return 'bg-purple-400/10';
  if (mimeType.includes('video')) return 'bg-red-400/10';
  if (mimeType.includes('audio')) return 'bg-pink-400/10';
  if (mimeType.includes('zip') || mimeType.includes('rar') || mimeType.includes('archive')) {
    return 'bg-yellow-600/10';
  }
  
  return 'bg-gray-400/10';
};

export default function FileItemCard({ item, viewMode, onClick, onStar, isSelected, onSelect }: FileItemCardProps) {
  const [showMenu, setShowMenu] = useState(false);

  if (viewMode === 'list') {
    return (
      <div 
        onClick={onClick}
        className={`flex items-center gap-4 px-4 py-3 hover:bg-white/5 cursor-pointer border-b border-white/5 ${
          isSelected ? 'bg-blue-600/20' : ''
        }`}
      >
        <button 
          onClick={(e) => { e.stopPropagation(); onSelect(); }}
          className={`w-5 h-5 rounded border-2 flex items-center justify-center transition-colors ${
            isSelected ? 'bg-blue-500 border-blue-500' : 'border-gray-500 hover:border-blue-400'
          }`}
        >
          {isSelected && <div className="w-2 h-2 bg-white rounded-sm" />}
        </button>
        
        <div className={`w-10 h-10 rounded-lg flex items-center justify-center ${getIconBg(item)}`}>
          {getIcon(item)}
        </div>
        
        <div className="flex-1 min-w-0">
          <div className="flex items-center gap-2">
            <span className="text-white truncate font-medium">{item.name}</span>
            {item.starred && <Star className="w-4 h-4 text-yellow-400 fill-yellow-400 flex-shrink-0" />}
            {item.shared && <span className="text-xs text-blue-400">Shared</span>}
          </div>
        </div>
        
        <span className="text-gray-400 text-sm w-32">{item.owner === 'me' ? 'me' : item.owner}</span>
        <span className="text-gray-400 text-sm w-32">{formatDate(item.modifiedAt)}</span>
        <span className="text-gray-400 text-sm w-24 text-right">{formatFileSize(item.size)}</span>
        
        <div className="relative">
          <button 
            onClick={(e) => { e.stopPropagation(); setShowMenu(!showMenu); }}
            className="p-1 text-gray-400 hover:text-white hover:bg-white/10 rounded"
          >
            <MoreVertical className="w-5 h-5" />
          </button>
          
          {showMenu && (
            <div className="absolute right-0 top-full mt-1 w-48 bg-[#2d2d2d] border border-white/10 rounded-lg shadow-xl z-10 py-1">
              <button className="w-full flex items-center gap-2 px-3 py-2 text-gray-300 hover:bg-white/10 text-sm">
                <Share2 className="w-4 h-4" /> Share
              </button>
              <button className="w-full flex items-center gap-2 px-3 py-2 text-gray-300 hover:bg-white/10 text-sm">
                <Download className="w-4 h-4" /> Download
              </button>
              <button className="w-full flex items-center gap-2 px-3 py-2 text-gray-300 hover:bg-white/10 text-sm">
                <Copy className="w-4 h-4" /> Make a copy
              </button>
              <button className="w-full flex items-center gap-2 px-3 py-2 text-gray-300 hover:bg-white/10 text-sm">
                <Pencil className="w-4 h-4" /> Rename
              </button>
              <hr className="my-1 border-white/10" />
              <button className="w-full flex items-center gap-2 px-3 py-2 text-red-400 hover:bg-white/10 text-sm">
                <Trash2 className="w-4 h-4" /> Remove
              </button>
            </div>
          )}
        </div>
      </div>
    );
  }

  return (
    <div 
      onClick={onClick}
      className={`group relative p-4 rounded-xl hover:bg-white/5 cursor-pointer transition-colors ${
        isSelected ? 'bg-blue-600/20 ring-2 ring-blue-500' : ''
      }`}
    >
      <button 
        onClick={(e) => { e.stopPropagation(); onSelect(); }}
        className={`absolute top-3 left-3 w-5 h-5 rounded border-2 flex items-center justify-center transition-all opacity-0 group-hover:opacity-100 ${
          isSelected ? 'opacity-100 bg-blue-500 border-blue-500' : 'border-gray-500 hover:border-blue-400 bg-[#1f1f1f]'
        }`}
      >
        {isSelected && <div className="w-2 h-2 bg-white rounded-sm" />}
      </button>
      
      <button 
        onClick={(e) => { e.stopPropagation(); onStar(); }}
        className={`absolute top-3 right-3 p-1 rounded transition-opacity opacity-0 group-hover:opacity-100 ${
          item.starred ? 'opacity-100' : ''
        }`}
      >
        <Star className={`w-5 h-5 ${item.starred ? 'text-yellow-400 fill-yellow-400' : 'text-gray-400'}`} />
      </button>
      
      <div className={`w-full aspect-square rounded-lg flex items-center justify-center mb-3 ${getIconBg(item)}`}>
        {getIcon(item)}
      </div>
      
      <div className="flex items-start justify-between gap-2">
        <div className="flex-1 min-w-0">
          <h3 className="text-white text-sm font-medium truncate">{item.name}</h3>
          <p className="text-gray-500 text-xs mt-0.5">
            {item.type === 'folder' ? 'Folder' : formatFileSize(item.size)}
          </p>
        </div>
      </div>
      
      {item.shared && (
        <div className="absolute bottom-12 left-4 text-xs text-blue-400">Shared</div>
      )}
    </div>
  );
}
