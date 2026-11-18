const BASE=(import.meta.env.VITE_API_URL as string|undefined)?.replace(/\/$/,'')??'';
let token='';
export function authorization():Record<string,string>{return token?{Authorization:`Bearer ${token}`}:{}}
export async function authCall(path:string,method='GET',body?:unknown){
 if(!BASE)throw new Error('Board access is unavailable until the API is configured.');
 const res=await fetch(`${BASE}/api/v1${path}`,{method,headers:{'Content-Type':'application/json',...authorization()},body:body===undefined?undefined:JSON.stringify(body),signal:AbortSignal.timeout(10000)});
 if(!res.ok){if(res.status===401)token='';const error=await res.json().catch(()=>({}));throw new Error(error.error||'The request could not be completed. Please retry.');}
 return res.status===204?undefined:res.json();
}
export const auth={
 async me(){if(!token)return null;return authCall('/auth/me')},
 async login(data:unknown){const result=await authCall('/auth/login','POST',data);token=result.token;return authCall('/auth/me')},
 logout(){token=''},
 activate:(data:unknown)=>authCall('/auth/activate','POST',data),
 members:()=>authCall('/admin/members'),
 invite:(data:unknown)=>authCall('/admin/members','POST',data),
 revoke:(email:string)=>authCall('/admin/members/'+encodeURIComponent(email),'DELETE')
};
