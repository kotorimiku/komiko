export interface ComicReadConfig {
  layoutMode: 'single' | 'double';
  direction: 'rightToLeft' | 'leftToRight';
  isSplit: boolean;
  matchSpreadPages: boolean;
}

export function getComicReadConfig(): ComicReadConfig {
  const readConfig = localStorage.getItem('comicReadConfig');
  if (readConfig) {
    return JSON.parse(readConfig) as ComicReadConfig;
  }
  return {
    layoutMode: 'single',
    direction: 'rightToLeft',
    isSplit: true,
    matchSpreadPages: false,
  };
}

export function setComicReadConfig(readConfig: ComicReadConfig) {
  localStorage.setItem('comicReadConfig', JSON.stringify(readConfig));
}