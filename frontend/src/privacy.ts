const CHOICE='aia-storage-choice-v1';
const keys=new Set(['aia-field-notes-v1','aia-field-notes-client']);
const memory=new Map<string,string>();
const originalGet=Storage.prototype.getItem, originalSet=Storage.prototype.setItem, originalRemove=Storage.prototype.removeItem;
let choice:string|null=null;
try { choice=originalGet.call(localStorage,CHOICE); } catch { /* blocked storage */ }
let allowed=choice==='allow';
Storage.prototype.getItem=function(key:string){return this===localStorage&&keys.has(key)&&!allowed?memory.get(key)??null:originalGet.call(this,key)};
Storage.prototype.setItem=function(key:string,value:string){if(this===localStorage&&keys.has(key)&&!allowed){memory.set(key,String(value));return}try{originalSet.call(this,key,value)}catch{if(keys.has(key))memory.set(key,String(value));else throw new Error('Storage is unavailable.')}};
Storage.prototype.removeItem=function(key:string){memory.delete(key);originalRemove.call(this,key)};
export const privacy={
 allowed:()=>allowed,
 choice:()=>choice,
 async set(on:boolean){
  if(!on){for(const key of keys){const value=localStorage.getItem(key);if(value)memory.set(key,value);originalRemove.call(localStorage,key)}
   const request=indexedDB.deleteDatabase('aia-field-notes-media');request.onerror=()=>console.warn('Local media deletion failed; use browser site-data settings.');request.onblocked=()=>console.warn('Close other site tabs to finish removing saved media.');
  }
  allowed=on;choice=on?'allow':'necessary';
  try{originalSet.call(localStorage,CHOICE,choice);if(on)for(const [k,v]of memory)originalSet.call(localStorage,k,v)}catch{/* current visit preference still applies */}
 }
};
