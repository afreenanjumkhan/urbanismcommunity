/* Readable source for the canonical prototype enhancement layer. */
(() => {
 const bridge=()=>window.fieldNotes;
 function node(tag,text,cls){const n=document.createElement(tag);if(text)n.textContent=text;if(cls)n.className=cls;return n}
 function button(text,action,cls){const b=node('button',text,cls);b.type='button';b.addEventListener('click',action);return b}
 function panel(title){const d=node('dialog',null,'fn-panel');const heading=node('h2',title);heading.id='fn-dialog-title';d.setAttribute('aria-labelledby',heading.id);const previous=document.activeElement;d.append(button('Close',()=>d.close(),'fn-close'),heading);d.addEventListener('close',()=>{d.remove();previous?.focus()});document.body.append(d);d.showModal();return d}
 function status(d,text){let p=d.querySelector('[role=alert]');if(!p){p=node('p');p.setAttribute('role','alert');d.append(p)}p.textContent=text}
 function field(form,label,type,name){const l=node('label',label),i=node('input');i.type=type;i.name=name;i.required=true;i.autocomplete=name==='email'?'email':name==='password'?'current-password':'off';l.append(i);form.append(l);return i}
 function privacy(){const d=panel('Privacy & storage');d.append(node('p','Public journal entries and session contributions are shared records. Please do not submit sensitive personal information. Board sign-in uses your email and a salted password hash; access tokens stay in memory and expire after one hour.'),node('p','Necessary storage remembers your privacy choice. Optional storage remembers your notes, bookmarks, anonymous voting identifier, and uploaded media on this device. Without it, activity stays in memory for this visit. We do not add advertising or analytics trackers.'),node('p','Accepting storage does not make public contributions private. Withdrawing storage permission removes this site’s saved local activity and media; it does not delete contributions already submitted to the shared record.'),node('p','For privacy, access, correction, or accessibility requests, contact urbanism@aiabrooklyn.org. Hosting providers may process operational logs. The committee must review retention and applicable legal obligations.'),button('Change storage choice',()=>{d.close();consent(true)}));}
 function consent(force=false){if(!bridge()?.privacy)return;if(!force&&bridge().privacy.choice())return;document.querySelector('.fn-consent')?.remove();const d=node('section',null,'fn-panel fn-consent');d.setAttribute('aria-label','Storage preferences');d.append(node('h2','Your notes. Your choice.'),node('p','May we remember your notes, bookmarks and uploaded media on this device? Necessary storage keeps your choice. No advertising or analytics cookies.'));const actions=node('div',null,'fn-actions');for(const [text,on] of [['Necessary only',false],['Allow local storage',true]])actions.append(button(text,async()=>{await bridge().privacy.set(on);d.remove()},on?'primary':''));actions.append(button('Privacy details',privacy));d.append(actions);document.body.append(d)}
 async function account(){const d=panel('Board access');if(!bridge()?.auth){d.append(node('p','Board access is available on the connected website.'));return}
 try{const user=await bridge().auth.me();if(user){d.append(node('p',`${user.email} · ${user.role==='super_admin'?'Super admin':'Board member'}`),button('Sign out',()=>{bridge().auth.logout();d.close()}));if(user.role==='super_admin')await members(d);return}}catch{}
 d.append(node('p','Sign in to create sessions. Public reading and conversations remain open.'));
 const form=node('form');field(form,'Email','email','email');field(form,'Password','password','password');const submit=node('button','Sign in','primary');submit.type='submit';form.append(submit);form.onsubmit=async e=>{e.preventDefault();submit.disabled=true;try{await bridge().auth.login(Object.fromEntries(new FormData(form)));d.close();account()}catch(e){status(d,e.message)}finally{submit.disabled=false}};d.append(form,button('Activate an invitation / first admin setup',()=>{d.close();activate()}));
 }
 function activate(){const d=panel('Activate your account');d.append(node('p','Use the private invitation from the super admin. First admin setup requires the deployment setup token.'));const f=node('form');field(f,'Email','email','email');field(f,'Invitation or setup token','password','token');const pass=field(f,'Choose a password (12–128 characters)','password','password');pass.minLength=12;pass.maxLength=128;pass.autocomplete='new-password';const b=node('button','Activate account','primary');f.append(b);f.onsubmit=async e=>{e.preventDefault();b.disabled=true;try{await bridge().auth.activate(Object.fromEntries(new FormData(f)));d.close();account()}catch(e){status(d,e.message)}finally{b.disabled=false}};d.append(f)}
 async function members(d){d.append(node('h3','Board membership'),node('p','Invite a board member, then share the one-time code privately. Codes expire in 24 hours. Revoking membership immediately blocks further authenticated requests.'));const f=node('form');field(f,'Board member email','email','email');const b=node('button','Create invitation','primary');f.append(b);f.onsubmit=async e=>{e.preventDefault();b.disabled=true;try{const result=await bridge().auth.invite(Object.fromEntries(new FormData(f)));let out=d.querySelector('output');if(!out){out=node('output');d.append(out)}out.textContent=`Invitation for ${new FormData(f).get('email')}: ${result.invitation}`;out.setAttribute('aria-live','polite')}catch(e){status(d,e.message)}finally{b.disabled=false}};d.append(f);const list=node('div');for(const u of await bridge().auth.members()){const row=node('div',null,'fn-member');row.append(node('span',`${u.email} · ${u.role} · ${u.active?'active':'inactive / invited'}`));if(u.role!=='super_admin')row.append(button('Revoke access',async()=>{try{await bridge().auth.revoke(u.email);row.textContent=`${u.email} · access revoked`}catch(e){status(d,e.message)}}));list.append(row)}d.append(list)}
 window.fnAccount=account;
 window.fnCreateSession=async(session,save)=>{if(!bridge()?.auth){account();return}try{await bridge().createSession(session);save()}catch(e){const d=panel('Session not saved');d.append(node('p',e.message),button('Board sign-in',()=>{d.close();account()}))}};
 window.fnRequireBoard=async(open)=>{try{const u=await bridge()?.auth.me();if(u&&['board','super_admin'].includes(u.role)){open();return}}catch{}account()};
 /* Home page (spec 003): the bundle mounts an empty #fn-home for #/home; this fills it from the
    shared document. Only real records are shown, and user text is set with textContent. */
 const SERIES=[['/about','About','Who we are, and why the committee keeps this record.'],['/journal','Journal','Field notes and essays from Brooklyn’s streets.'],['/sessions','Sessions','Upcoming conversations. Come and add your voice.'],['/archives','Archives','Past sessions and what neighbors said.'],['/records','Records','The shared record, folder by folder.'],['/contact','Contact','Reach the committee and join the community channels.']];
 const readDoc=()=>{try{return JSON.parse(localStorage.getItem('aia-field-notes-v1')||'null')}catch{return null}};
 const calm=()=>matchMedia('(prefers-reduced-motion: reduce)').matches;
 const asDate=d=>new Date(`${d}T12:00:00`);
 function anchor(text,href,cls){const a=node('a',text,cls);a.href=href;return a}
 function eyebrow(text){const d=node('div',null,'eyebrow');d.append(node('span'),text);return d}
 function renderHome(el){
  el.dataset.ready='1';
  const doc=el.fnState||readDoc(),sessions=doc?.sessions||[],byId=Object.fromEntries(sessions.map(s=>[s.id,s]));
  const notes=(doc?.entries||[]).filter(e=>e&&e.text).sort((a,b)=>String(b.created).localeCompare(String(a.created)));
  // Hero copy
  const hero=node('section',null,'fn-home-hero');hero.setAttribute('aria-labelledby','fn-home-title');
  const copy=node('div',null,'fn-home-copy'),h1=node('h1');h1.id='fn-home-title';h1.append('Brooklyn, ',node('span','block by block.','fn-mark'));
  const actions=node('div',null,'fn-home-actions');actions.append(anchor('Read the journal','#/journal','button red'),anchor('Find a session','#/sessions','text-link'));
  copy.append(eyebrow('AIA BROOKLYN / URBANISM COMMITTEE'),h1,node('p','Field Notes is an open record of what neighbors notice on their streets: the corners, crossings and waterfronts we share. Every note joins the conversations that shape what comes next for Brooklyn.'),actions);
  // Interactive block map
  const city=node('div',null,'fn-city'),head=node('div',null,'fn-city-head');head.append(eyebrow('A MAP OF SMALL ENCOUNTERS'),node('small','Move across the blocks. Tap one to read a field note.'));
  const grid=node('div',null,'fn-city-grid');grid.setAttribute('aria-hidden','true');
  const COLS=12,ROWS=8,blocks=[];
  for(let i=0;i<COLS*ROWS;i++){const b=node('i'),r=Math.floor(i/COLS),c=i%COLS;
   if(c===4||r===3)b.className='street';else if(r>=5&&c>=8&&c<=10)b.className='park';else b.style.setProperty('--b',`var(--red-${1+((i*7+r*3)%4)})`);
   blocks.push({b,x:(c+.5)/COLS,y:(r+.5)/ROWS,lot:!b.className});grid.append(b)}
  const lots=blocks.filter(k=>k.lot);
  const note=node('figure',null,'fn-note'),body=node('div',null,'fn-note-body'),controls=node('div',null,'fn-note-actions');body.setAttribute('aria-live','polite');
  const from=anchor('','#/sessions','text-link');let current=-1,pinned=null;
  function show(i,block){
   if(!notes.length){body.replaceChildren(node('p','No field notes yet. Join a session and leave the first one.'));from.textContent='See sessions';return}
   current=(i+notes.length)%notes.length;const n=notes[current],s=byId[n.sessionId];
   pinned?.classList.remove('pin');pinned=block||lots[(current*7)%lots.length].b;pinned.classList.add('pin');
   body.replaceChildren(eyebrow(`FIELD NOTE ${String(current+1).padStart(2,'0')} / ${String(notes.length).padStart(2,'0')}`),node('blockquote',`“${n.text}”`),node('figcaption',[n.name||'A neighbor',n.location,n.hood].filter(Boolean).join(' · ')));
   from.hidden=!s;if(s){from.textContent=`From “${s.title}”`;from.href=`#/session/${s.code}`}
   if(!calm()){body.classList.remove('fn-in');void body.offsetWidth;body.classList.add('fn-in')}}
  controls.append(button('Show another field note',()=>show(current+1),'text-link'),from);note.append(body,controls);
  let raf=0,px=-9,py=-9,hover=false;
  const paint=()=>{raf=0;for(const k of lots){const d=Math.hypot((k.x-px)*COLS,(k.y-py)*ROWS);k.b.style.setProperty('--lift',Math.max(0,1-d/3).toFixed(2))}};
  const aim=(x,y)=>{px=x;py=y;if(!raf)raf=requestAnimationFrame(paint)};
  grid.addEventListener('pointermove',e=>{const r=grid.getBoundingClientRect();hover=true;aim((e.clientX-r.left)/r.width,(e.clientY-r.top)/r.height)});
  grid.addEventListener('pointerleave',()=>{hover=false;aim(-9,-9)});
  grid.addEventListener('click',e=>{const at=lots.findIndex(k=>k.b===e.target);if(at>=0)show(at,lots[at].b)});
  if(!calm()){let t=0;const tick=setInterval(()=>{if(!el.isConnected){clearInterval(tick);return}if(hover||document.hidden)return;t+=.035;aim(.5+.45*Math.sin(t),.5+.4*Math.sin(t*1.7))},60)}
  city.append(head,grid,note);hero.append(copy,city);el.append(hero);show(0);
  if(!doc){el.append(SERIES_SECTION());return}
  // The record so far
  const facts=[[sessions.length,'Sessions'],[notes.length,'Field notes'],[new Set(notes.map(n=>n.hood).filter(Boolean)).size,'Neighborhoods heard'],[(doc.replies||[]).length,'Replies']];
  const stats=node('section',null,'fn-home-stats');stats.setAttribute('aria-label','The record so far');
  for(const [n,label] of facts){const d=node('div'),b=node('b',String(n));b.dataset.to=String(n);d.append(b,node('span',label));stats.append(d)}
  if(!calm()&&'IntersectionObserver' in window){stats.querySelectorAll('b').forEach(b=>b.textContent='0');const io=new IntersectionObserver(es=>{if(!es.some(e=>e.isIntersecting))return;io.disconnect();const t0=performance.now(),step=now=>{const p=Math.min(1,(now-t0)/1100),ease=1-Math.pow(1-p,3);stats.querySelectorAll('b').forEach(b=>b.textContent=String(Math.round(ease*Number(b.dataset.to))));if(p<1)requestAnimationFrame(step)};requestAnimationFrame(step)});io.observe(stats)}
  el.append(stats);
  // Next session, or the most recent one
  const today=new Date().toLocaleDateString('en-CA',{timeZone:'America/New_York'});
  const next=sessions.filter(s=>s.date>=today).sort((a,b)=>a.date.localeCompare(b.date))[0],pick=next||[...sessions].sort((a,b)=>b.date.localeCompare(a.date))[0];
  if(pick){const sec=node('section',null,'fn-home-next'),d=asDate(pick.date),cal=node('div',null,'fn-next-date'),info=node('div'),h=node('h2',pick.title);h.id='fn-next-title';sec.setAttribute('aria-labelledby',h.id);
   cal.append(node('span',d.toLocaleDateString('en-US',{month:'short'}).toUpperCase()),node('b',String(d.getDate())),node('span',String(d.getFullYear())));
   info.append(eyebrow(next?'NEXT ON THE CALENDAR':'MOST RECENT SESSION'),h,node('p',[d.toLocaleDateString('en-US',{weekday:'long',month:'long',day:'numeric'}),pick.place,pick.neighborhood].filter(Boolean).join(' · ')));
   sec.append(cal,info,anchor(next?'Join this session':'Read the session',`#/session/${pick.code}`,'button red'));el.append(sec)}
  el.append(SERIES_SECTION());
 }
 function SERIES_SECTION(){const wrap=node('section',null,'fn-home-folders-wrap'),h=node('h2','Open a folder.'),list=node('div',null,'fn-home-folders');h.id='fn-folders-title';wrap.setAttribute('aria-labelledby',h.id);
  SERIES.forEach(([href,title,text],i)=>{const a=anchor(null,`#${href}`,'fn-folder'),arrow=node('span','↗','fn-folder-arrow');arrow.setAttribute('aria-hidden','true');a.style.setProperty('--tab',`var(--red-${i+2})`);a.append(node('span',String(i+2).padStart(2,'0'),'fn-folder-no'),node('h3',title),node('p',text),arrow);list.append(a)});
  wrap.append(eyebrow('THE SERIES'),h,list);return wrap}
 const observed=new WeakMap();let menu;
 function enhance(){
  const home=document.getElementById('fn-home');if(home&&!home.dataset.ready)renderHome(home);
  const nav=document.querySelector('.main-nav');if(nav&&!document.querySelector('.fn-menu')){nav.id='fn-navigation';menu=button('☰ Menu',()=>{const open=nav.classList.toggle('fn-open');menu.setAttribute('aria-expanded',String(open));if(open)nav.querySelector('a')?.focus()},'fn-menu');menu.setAttribute('aria-controls',nav.id);menu.setAttribute('aria-expanded','false');nav.before(menu);nav.addEventListener('click',e=>{if(e.target.closest('a'))closeMenu()})}
  document.querySelectorAll('img').forEach(img=>{if(observed.get(img)===img.src)return;observed.set(img,img.src);img.classList.remove('fn-failed');const done=()=>{img.classList.remove('fn-loading');img.removeAttribute('aria-busy');if(!img.naturalWidth){img.classList.add('fn-failed');img.alt=`Image unavailable. ${img.alt.replace(/^Image unavailable\. /,'')}`}};if(img.complete){done();return}img.classList.add('fn-loading');img.setAttribute('aria-busy','true');img.addEventListener('load',done,{once:true});img.addEventListener('error',done,{once:true});img.decoding='async'});
  if(document.querySelector('.site-footer')&&!document.querySelector('.fn-tools')){const bar=node('div',null,'fn-tools');bar.append(button('Board sign-in / administration',account),button('Privacy & storage',privacy),button('Storage preferences',()=>consent(true)));document.querySelector('.site-footer').append(bar);consent()}
 }
 function closeMenu(){document.querySelector('.main-nav')?.classList.remove('fn-open');menu?.setAttribute('aria-expanded','false')}
 document.addEventListener('keydown',e=>{if(e.key==='Escape'&&menu?.getAttribute('aria-expanded')==='true'){closeMenu();menu.focus()}});
 window.addEventListener('hashchange',()=>{closeMenu();setTimeout(()=>document.getElementById('main')?.focus(),0)});
 new MutationObserver(enhance).observe(document.getElementById('root'),{childList:true,subtree:true});enhance();
})();
