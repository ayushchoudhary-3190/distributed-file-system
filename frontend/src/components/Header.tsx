'use client';

import { Search, Grid, List, HelpCircle, Settings, Bell, User } from 'lucide-react';
import { useState } from 'react';
import { ViewMode } from '@/types';

interface HeaderProps {
  searchQuery: string;
  onSearchChange: (query: string) => void;
  viewMode: ViewMode;
  onViewModeChange: (mode: ViewMode) => void;
}

export default function Header({ searchQuery, onSearchChange, viewMode, onViewModeChange }: HeaderProps) {
  return (
    <header className="h-16 bg-[#1f1f1f] border-b border-white/10 flex items-center justify-between px-4">
      <div className="flex items-center gap-4 flex-1">
        <div className="flex items-center gap-2 text-gray-400">
          <span className="text-sm">My Drive</span>
        </div>
      </div>

      <div className="flex-1 max-w-2xl">
        <div className="relative">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-400" />
          <input
            type="text"
            value={searchQuery}
            onChange={(e) => onSearchChange(e.target.value)}
            placeholder="Search in Drive"
            className="w-full pl-10 pr-4 py-2 bg-[#2d2d2d] border border-transparent focus:border-blue-500 rounded-full text-white placeholder-gray-500 outline-none transition-colors"
          />
        </div>
      </div>

      <div className="flex items-center gap-2 ml-4">
        <button className="p-2 text-gray-400 hover:text-white hover:bg-white/10 rounded-full transition-colors">
          <HelpCircle className="w-5 h-5" />
        </button>
        <button className="p-2 text-gray-400 hover:text-white hover:bg-white/10 rounded-full transition-colors">
          <Settings className="w-5 h-5" />
        </button>
        <button className="p-2 text-gray-400 hover:text-white hover:bg-white/10 rounded-full transition-colors">
          <Bell className="w-5 h-5" />
        </button>
        
        <div className="flex items-center gap-2 ml-2 pl-2 border-l border-white/10">
          <button 
            onClick={() => onViewModeChange('grid')}
            className={`p-2 rounded-full transition-colors ${
              viewMode === 'grid' ? 'text-blue-400 bg-blue-400/20' : 'text-gray-400 hover:text-white hover:bg-white/10'
            }`}
          >
            <Grid className="w-5 h-5" />
          </button>
          <button 
            onClick={() => onViewModeChange('list')}
            className={`p-2 rounded-full transition-colors ${
              viewMode === 'list' ? 'text-blue-400 bg-blue-400/20' : 'text-gray-400 hover:text-white hover:bg-white/10'
            }`}
          >
            <List className="w-5 h-5" />
          </button>
        </div>

        <button className="ml-2 w-9 h-9 bg-gradient-to-br from-purple-500 to-pink-500 rounded-full flex items-center justify-center text-white font-medium">
          <User className="w-5 h-5" />
        </button>
      </div>
    </header>
  );
}
