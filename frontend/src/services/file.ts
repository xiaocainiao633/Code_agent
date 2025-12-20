import { fileAPI } from './api'
import type { FileInfo } from '@/types/task'

// 文件上传配置
export interface UploadConfig {
  maxSize: number // 最大文件大小 (字节)
  allowedTypes: string[] // 允许的文件类型
  chunkSize: number // 分块大小 (字节)
}

// 默认上传配置
export const DEFAULT_UPLOAD_CONFIG: UploadConfig = {
  maxSize: 10 * 1024 * 1024, // 10MB
  allowedTypes: ['.py', '.js', '.ts', '.java', '.cpp', '.c', '.h', '.txt', '.md'],
  chunkSize: 1024 * 1024, // 1MB
}

// 文件服务类
export class FileService {
  private uploadConfig: UploadConfig

  constructor(config: UploadConfig = DEFAULT_UPLOAD_CONFIG) {
    this.uploadConfig = config
  }

  // 验证文件
  validateFile(file: File): { valid: boolean; error?: string } {
    // 检查文件大小
    if (file.size > this.uploadConfig.maxSize) {
      return {
        valid: false,
        error: `文件大小超过限制 (${this.formatFileSize(this.uploadConfig.maxSize)})`
      }
    }

    // 检查文件类型
    const fileExtension = '.' + file.name.split('.').pop()?.toLowerCase()
    if (!this.uploadConfig.allowedTypes.includes(fileExtension)) {
      return {
        valid: false,
        error: `不支持的文件类型. 支持的类型: ${this.uploadConfig.allowedTypes.join(', ')}`
      }
    }

    return { valid: true }
  }

  // 格式化文件大小
  formatFileSize(bytes: number): string {
    if (bytes === 0) return '0 Bytes'
    
    const k = 1024
    const sizes = ['Bytes', 'KB', 'MB', 'GB']
    const i = Math.floor(Math.log(bytes) / Math.log(k))
    
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
  }

  // 读取文件内容
  readFileContent(file: File): Promise<string> {
    return new Promise((resolve, reject) => {
      const reader = new FileReader()
      
      reader.onload = (e) => {
        const content = e.target?.result as string
        resolve(content)
      }
      
      reader.onerror = () => {
        reject(new Error('读取文件失败'))
      }
      
      reader.readAsText(file)
    })
  }

  // 上传文件
  async uploadFile(file: File): Promise<FileInfo> {
    // 验证文件
    const validation = this.validateFile(file)
    if (!validation.valid) {
      throw new Error(validation.error)
    }

    try {
      console.log(`Uploading file: ${file.name} (${this.formatFileSize(file.size)})`)
      const result = await fileAPI.uploadFile(file)
      console.log(`File uploaded successfully: ${result.filename}`)
      return result
    } catch (error) {
      console.error('File upload failed:', error)
      throw new Error(`文件上传失败: ${error instanceof Error ? error.message : '未知错误'}`)
    }
  }

  // 批量上传文件
  async uploadFiles(files: File[]): Promise<FileInfo[]> {
    const results: FileInfo[] = []
    const errors: string[] = []

    for (const file of files) {
      try {
        const result = await this.uploadFile(file)
        results.push(result)
      } catch (error) {
        errors.push(`${file.name}: ${error instanceof Error ? error.message : '上传失败'}`)
      }
    }

    if (errors.length > 0) {
      console.warn(`Some files failed to upload: ${errors.join(', ')}`)
    }

    return results
  }

  // 获取文件列表
  async getFiles(): Promise<FileInfo[]> {
    try {
      const files = await fileAPI.getFiles()
      console.log(`Retrieved ${files.length} files`)
      return files
    } catch (error) {
      console.error('Failed to get files:', error)
      throw new Error(`获取文件列表失败: ${error instanceof Error ? error.message : '未知错误'}`)
    }
  }

  // 删除文件
  async deleteFile(fileId: string): Promise<void> {
    try {
      await fileAPI.deleteFile(fileId)
      console.log(`File deleted: ${fileId}`)
    } catch (error) {
      console.error('Failed to delete file:', error)
      throw new Error(`删除文件失败: ${error instanceof Error ? error.message : '未知错误'}`)
    }
  }

  // 批量删除文件
  async deleteFiles(fileIds: string[]): Promise<void> {
    const errors: string[] = []

    for (const fileId of fileIds) {
      try {
        await this.deleteFile(fileId)
      } catch (error) {
        errors.push(`${fileId}: ${error instanceof Error ? error.message : '删除失败'}`)
      }
    }

    if (errors.length > 0) {
      throw new Error(`部分文件删除失败: ${errors.join(', ')}`)
    }
  }

  // 批量处理文件
  async batchProcessFiles(fileIds: string[], operation: string, params: any = {}) {
    try {
      console.log(`Batch processing ${fileIds.length} files with operation: ${operation}`)
      const result = await fileAPI.batchProcessFiles(fileIds, operation, params)
      console.log(`Batch processing completed`)
      return result
    } catch (error) {
      console.error('Batch processing failed:', error)
      throw new Error(`批量处理失败: ${error instanceof Error ? error.message : '未知错误'}`)
    }
  }

  // 下载文件
  downloadFile(file: FileInfo): void {
    // 创建下载链接
    const link = document.createElement('a')
    link.href = `/api/v1/files/${file.id}/download`
    link.download = file.filename
    link.click()
  }

  // 获取文件图标
  getFileIcon(fileName: string): string {
    const extension = fileName.split('.').pop()?.toLowerCase()
    
    const iconMap: Record<string, string> = {
      'py': '🔧',
      'js': '📜',
      'ts': '📘',
      'java': '☕',
      'cpp': '⚙️',
      'c': '⚙️',
      'h': '📋',
      'txt': '📄',
      'md': '📝',
    }

    return iconMap[extension || ''] || '📎'
  }

  // 获取文件类型描述
  getFileTypeDescription(fileName: string): string {
    const extension = fileName.split('.').pop()?.toLowerCase()
    
    const typeMap: Record<string, string> = {
      'py': 'Python 文件',
      'js': 'JavaScript 文件',
      'ts': 'TypeScript 文件',
      'java': 'Java 文件',
      'cpp': 'C++ 文件',
      'c': 'C 文件',
      'h': '头文件',
      'txt': '文本文件',
      'md': 'Markdown 文件',
    }

    return typeMap[extension || ''] || '未知文件类型'
  }
}

// 创建全局文件服务实例
export const fileService = new FileService()

export default fileService