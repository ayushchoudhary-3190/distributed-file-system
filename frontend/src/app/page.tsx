'use client';

import { useState } from 'react';
import Sidebar from '@/components/Sidebar';
import Header from '@/components/Header';
import FileGrid from '@/components/FileGrid';
import UploadModal from '@/components/UploadModal';
import { mockFiles } from '@/data/files';
import { ViewMode } from '@/types';
import { Plus, ArrowUp, ChevronRight, ChevronDown, Square, CheckSquare } from 'lucide-react';

export default function Home() {
  const [activeSection, setActiveSection] = useState('my-drive');
  const [searchQuery, setSearchQuery] = useState('');
  const [viewMode, setViewMode] = useState<ViewMode>('grid');
  const [showUploadModal, setShowUploadModal] = useState(false);
  const [sortField, setSortField] = useState<'name' | 'modified'>('name');
  const [sortAsc, setSortAsc] = useState(true);

  const filteredFiles = mockFiles
    .filter(file => {
      if (activeSection === 'starred') return file.starred;
      if (activeSection === 'shared') return file.shared;
      return true;
    })
    .filter(file => 
      file.name.toLowerCase().includes(searchQuery.toLowerCase())
    )
    .sort((a, b) => {
      if (a.type !== b.type) return a.type === 'folder' ? -1 : 1;
      
      let comparison = 0;
      if (sortField === 'name') {
        comparison = a.name.localeCompare(b.name);
      } else if (sortField === 'modified') {
        comparison = a.modifiedAt.getTime() - b.modifiedAt.getTime();
      }
      
      return sortAsc ? comparison : -comparison;
    });

  const getSectionTitle = () => {
    switch (activeSection) {
      case 'my-drive': return 'My Drive';
      case 'shared': return 'Shared with me';
      case 'recent': return 'Recent';
      case 'starred': return 'Starred';
      case 'trash': return 'Trash';
      default: return 'My Drive';
    }
  };

  return (
    <div className="flex h-screen bg-[#1f1f1f]">
      <Sidebar activeSection={activeSection} onSectionChange={setActiveSection} />
      
      <div className="flex-1 flex flex-col min-w-0">
        <Header 
          searchQuery={searchQuery}
          onSearchChange={setSearchQuery}
          viewMode={viewMode}
          onViewModeChange={setViewMode}
        />
        
        <div className="px-6 py-4 flex items-center justify-between">
          <div className="flex items-center gap-2">
            <h1 className="text-xl font-medium text-white">{getSectionTitle()}</h1>
            <ChevronRight className="w-4 h-4 text-gray-500" />
          </div>
          
          <div className="flex items-center gap-2">
            <button 
              onClick={() => setSortAsc(!sortAsc)}
              className="flex items-center gap-1 px-3 py-1.5 text-sm text-gray-400 hover:text-white hover:bg-white/10 rounded-lg transition-colors"
            >
              <span>{sortAsc ? 'Ascending' : 'Descending'}</span>
            </button>
            
            <button 
              onClick={() => setShowUploadModal(true)}
              className="flex items-center gap-2 px-4 py-2 bg-blue-600 hover:bg-blue--full transition700 text-white rounded-colors"
            >
              <Plus className="w-5 h-5" />
              <span>New</span>
            </button>
          </div>
        </div>
        
        <FileGrid files={filteredFiles} viewMode={viewMode} />
      </div>
      
      <UploadModal isOpen={showUploadModal} onClose={() => setShowUploadModal(false)} />
    </div>
  );
}
