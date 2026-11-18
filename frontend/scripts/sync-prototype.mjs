import {readFileSync,writeFileSync} from 'node:fs';
import {createHash} from 'node:crypto';
const root=new URL('../../',import.meta.url);
const read=p=>readFileSync(new URL(p,root),'utf8');
const write=(p,s)=>writeFileSync(new URL(p,root),s);
let html=read('prototype/index.html');
// Embed the vector mark so the canonical prototype stays self-contained.
const logo=read('frontend/public/media/brand/aia-brooklyn.svg');
const logoURI='data:image/svg+xml;base64,'+Buffer.from(logo).toString('base64');
const oldLockup='(0,b.jsxs)("span",{className:"chapter-lockup",children:[(0,b.jsx)("span",{className:"chapter-monogram","aria-hidden":"true",children:"BK"}),(0,b.jsx)("span",{className:"chapter-seal",children:"AIA"}),(0,b.jsx)("span",{className:"chapter-name",children:"Brooklyn"})]})';
const newLockup=`(0,b.jsx)("img",{className:"chapter-vector-logo",src:${JSON.stringify(logoURI)},alt:"AIA Brooklyn",width:1305,height:256})`;
if(html.includes(oldLockup))html=html.replace(oldLockup,newLockup);
else if(html.includes('className:"chapter-vector-logo"'))html=html.replace(/\(0,b\.jsx\)\("img",\{className:"chapter-vector-logo",src:"[^"]+",alt:"AIA Brooklyn",width:1305,height:256\}\)/,newLockup);
else throw new Error('Cannot find canonical header logo; review prototype integration.');
// Spec 003: home route and series navigation. Each swap is idempotent and fails loudly
// if a new prototype export no longer contains its anchor.
const swap=(from,to,all=false)=>{if(html.includes(to))return;if(!html.includes(from))throw new Error('Spec 003 anchor missing: '+from.slice(0,80));html=all?html.replaceAll(from,to):html.replace(from,to)};
swap('[["/journal","Journal"],["/sessions","Sessions"],["/archives","Archives"],["/about","About"],["/records","Records"],["/contact","Contact"]]','[["/home","Home"],["/about","About"],["/journal","Journal"],["/sessions","Sessions"],["/archives","Archives"],["/records","Records"],["/contact","Contact"]]');
swap('r==="/archives"&&(0,b.jsx)("sup",{children:v.length.toString().padStart(2,"0")})','(0,b.jsx)("sup",{className:"folder-number","aria-hidden":"true",children:String(cA+1).padStart(2,"0")})');
swap('location.hash.slice(1)||"/journal"','location.hash.slice(1)||"/home"',true);
html=html.replace('f==="/home"&&(0,b.jsx)("div",{id:"fn-home",className:"container home-page"}),',''); // drop the first draft of the mount
// The ref hands the live app state (A) to the enhancement layer; React renders no children here.
swap('f==="/about"&&(0,b.jsx)(tv,{})','f==="/home"&&(0,b.jsx)("div",{id:"fn-home",className:"container home-page",ref:n=>{if(n)n.fnState=A}}),f==="/about"&&(0,b.jsx)(tv,{})');
swap('!["/journal","/sessions","/archives","/about","/contact","/records"].includes(f)','!["/home","/journal","/sessions","/archives","/about","/contact","/records"].includes(f)');
swap('{"/sessions":"Sessions",','{"/home":"Home","/sessions":"Sessions",');
swap('className:"brand chapter-brand",href:"#/journal"','className:"brand chapter-brand",href:"#/home"');
// Embed open-license web fonts from prototype/fonts/ so the prototype stays self-contained.
const fontURI=f=>'data:font/woff2;base64,'+readFileSync(new URL('prototype/fonts/'+f,root)).toString('base64');
const css=read('prototype/enhancements.css').replace(/url\("fonts\/([a-z0-9-]+\.woff2)"\)/g,(m,f)=>`url("${fontURI(f)}")`),js=read('prototype/enhancements.js');
html=html.replace(/\/\* BEGIN ACCESSIBLE REVISION \*\/[\s\S]*?\/\* END ACCESSIBLE REVISION \*\//g,'');
html=html.replace('</style>',`/* BEGIN ACCESSIBLE REVISION */\n${css}\n/* END ACCESSIBLE REVISION */</style>`);
html=html.replace('</script>',`/* BEGIN ACCESSIBLE REVISION */\n${js}\n/* END ACCESSIBLE REVISION */</script>`);
html=html.replace('onNew:()=>s(!0)','onNew:()=>window.fnRequireBoard(()=>s(!0))');
html=html.replace('function Uu(){return new Promise(', 'function Uu(){if(window.fieldNotes&&!window.fieldNotes.privacy.allowed())return Promise.reject(new Error("Allow local storage in Storage preferences before saving media."));return new Promise(');
html=html.replace('onSave:r=>{k(R=>({...R,sessions:[...R.sessions,r]}))&&(s(!1),nJ(`/session/${r.code}`),E(`Session created. Your code is ${r.code}.`))}', 'onSave:r=>{window.fnCreateSession(r,()=>{k(R=>({...R,sessions:[...R.sessions,r]}))&&(s(!1),nJ(`/session/${r.code}`),E(`Session created. Your code is ${r.code}.`))})}');
html=html.replace('Created sessions are saved in this browser.','Session creation requires board access and a successful server response.');
html=html.replaceAll('Your activity is saved in this browser.','Optional storage remembers activity on this device.');
html=html.replaceAll('Saved in this browser, kept in your notebook.','Kept for this visit; allow storage to remember it.');
html=html.replace('Links will appear here once the committee\'s contact details are confirmed.','Contact the committee using the channels above.');
html=html.replace('Links will appear here once the committee\\u2019s contact details are confirmed.','Contact the committee using the channels above.');
write('prototype/index.html',html);
for(const [tag,file]of [['style','app.css'],['script','app.js']])write('frontend/public/legacy/'+file,html.match(new RegExp(`<${tag}\\b[^>]*>([\\s\\S]*?)</${tag}>`))[1]);
const version=createHash('sha256').update(html).digest('hex').slice(0,16);
write('frontend/src/ui-version.ts',`// Generated from the canonical prototype.\nexport const UI_VERSION = '${version}';\n`);
write('frontend/index.html',read('frontend/index.html').replace(/\/legacy\/app\.css(?:\?v=[a-z0-9]+)?/,`/legacy/app.css?v=${version}`));
console.log('Revised prototype and regenerated legacy assets.');
