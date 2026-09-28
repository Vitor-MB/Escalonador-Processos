const COLORS=[['--p1','--p1d'],['--p2','--p2d'],['--p3','--p3d'],['--p4','--p4d'],['--p5','--p5d'],['--p6','--p6d']];

let rows=0;
function addRow(pre){
  rows++;
  const id='P'+rows;
  const tb=document.getElementById('procBody');
  const tr=document.createElement('tr');
  const [fill,dark]=COLORS[(rows-1)%COLORS.length];
  tr.innerHTML=`<td><span class="swatch" style="background:var(${fill});border:1px solid var(${dark})"></span></td>
    <td><input value="${id}" data-f="name"></td>
    <td><input type="number" min="0" value="${pre?pre.a:0}" data-f="a" style="width:60px"></td>
    <td><input type="number" min="1" value="${pre?pre.b:4}" data-f="b" style="width:60px"></td>
    <td><input type="number" min="0" value="${pre?pre.p:1}" data-f="p" style="width:60px"></td>
    <td><button class="rowbtn" onclick="this.closest('tr').remove()">×</button></td>`;
  tb.appendChild(tr);
}


document.getElementById('algo').addEventListener('change',syncAlgoUI);

function syncAlgoUI(){
  const a=document.getElementById('algo').value;
  document.getElementById('quantumBox').classList.toggle('show', a==='rr'||a==='rr-prio-aging');
  document.getElementById('agingBox').classList.toggle('show', a==='rr-prio-aging');
}

// Adiciona alguns Elementos de exemplo para os processos
[{a:0,b:5,p:2},{a:1,b:3,p:1},{a:2,b:8,p:3},{a:3,b:6,p:2},{a:4,b:2,p:1}].forEach(addRow);
syncAlgoUI();

function readProcs(){
  return [...document.querySelectorAll('#procBody tr')].map((tr,i)=>{
    const g=f=>tr.querySelector(`[data-f="${f}"]`).value;
    return {name:g('name')||('P'+(i+1)), a:+g('a'), b:+g('b'), p:+g('p'), color:COLORS[i%COLORS.length]};
  });
}

//Simula os ticks de execução 
function simulate(procs, algo, quantum, aging){
    const N=procs.length;
    const st=procs.map(p=>({name:p.name, arrival:p.a, burst:p.b, prio:p.p, remaining:p.b, start:null, finish:null}));
    let t = Math.min(...st.map(p=>p.arrival));
    const intervals=[];
    let lastRunning=null, segStart=null;
    let contextSwitches = 0;

    function pushTick(name){
        if(name===lastRunning){ /* mesmo tick (não faz nada) */ }
        
        else{
            if(lastRunning!==null) {
                intervals.push({name:lastRunning, start:segStart, end:t});
                contextSwitches++;
            }
            lastRunning=name; segStart=t;
        }
    }

    function closeSeg(){ 
        if(lastRunning!==null){ 
            intervals.push({name:lastRunning, start:segStart, end:t}); 
            lastRunning=null; 
        } 
    }   

    const arrived = ()=>st.filter(p=>p.arrival<=t && p.remaining>0);

    if(algo==='fcfs' || algo==='sjf' || algo==='prio-np'){
        let committed=null;
        let guard=0;
        while(st.some(p=>p.remaining>0) && guard<100000){
            guard++;
            const av=arrived();

            if(!committed || committed.remaining<=0){
                if(av.length===0){ 
                    t++; 
                    continue; 
                }

                let pick;
                if(algo==='fcfs') 
                    pick=av.reduce((a,b)=> a.arrival<b.arrival?a: b.arrival<a.arrival?b: (a.name<b.name?a:b));
                else if(algo==='sjf') 
                    pick=av.reduce((a,b)=> a.burst<b.burst?a: b.burst<a.burst?b: (a.arrival<b.arrival?a:b));
                else 
                    pick=av.reduce((a,b)=> a.prio<b.prio?a: b.prio<a.prio?b: (a.arrival<b.arrival?a:b));

                committed=pick;
            }
            if(committed.start===null) 
                committed.start=t;
            pushTick(committed.name);
            committed.remaining--;
            t++;
            if(committed.remaining<=0){ 
                committed.finish=t; 
                committed=null; 
            }
        }
        
        closeSeg();
    }
    else if(algo==='srtf' || algo==='prio-p'){
        let guard=0;
        while(st.some(p=>p.remaining>0) && guard<100000){
            guard++;
            const av=arrived();
            if(av.length===0){ 
                t++; 
                continue; 
            }

            let pick;
            if(algo==='srtf') 
                pick=av.reduce((a,b)=> a.remaining<b.remaining?a: b.remaining<a.remaining?b: (a.arrival<b.arrival?a:b));
            else 
                pick=av.reduce((a,b)=> a.prio<b.prio?a: b.prio<a.prio?b: (a.arrival<b.arrival?a:b));
            
            if(pick.start===null) 
                pick.start=t;
            pushTick(pick.name);
            pick.remaining--;
            t++;
            if(pick.remaining<=0) 
                pick.finish=t;
        }
        closeSeg();
    }
    else if(algo==='rr'){
        const q=[]; const queued=new Set();
        st.sort((a,b)=>a.arrival-b.arrival);
        let i=0, guard=0, run=null, runLeft=0;
        while(st.some(p=>p.remaining>0) && guard<200000){
            guard++;
            while(i<st.length && st[i].arrival<=t){ 
                if(!queued.has(st[i].name) && st[i].remaining>0 && st[i]!==run){ 
                    q.push(st[i]); 
                    queued.add(st[i].name);
                } 
                i++; 
            }

            if(!run){
                if(q.length===0){ 
                    t++; 
                    continue; 
                }
                run=q.shift(); 
                queued.delete(run.name); 
                runLeft=Math.min(quantum, run.remaining);
            }
            if(run.start===null) 
                run.start=t;

            pushTick(run.name);
            run.remaining--; 
            runLeft--; 
            t++;

            // Enfileira novos processos que chegaram durante a execução do processo atual 
            while(i<st.length && st[i].arrival<=t){ 
                if(!queued.has(st[i].name) && st[i].remaining>0 && st[i]!==run){ 
                    q.push(st[i]); 
                    queued.add(st[i].name);
                } 
                i++; 
            }
            if(run.remaining<=0){ 
                run.finish=t; 
                run=null; 
            }
            else if(runLeft<=0){ 
                q.push(run); 
                queued.add(run.name); 
                run=null; 
            }
        }
        closeSeg();
    }

    else if(algo==='rr-prio-aging'){
        const q=[]; 
        const queued=new Set();
        st.sort((a,b)=>a.arrival-b.arrival);
        const waitSince={}; 
        st.forEach(p=>waitSince[p.name]=null);
        let i=0, guard=0, run=null, runLeft=0;
        
        function effPrio(p){
            if(waitSince[p.name]===null) 
                return p.prio;

            const waited=t-waitSince[p.name];
            return Math.max(0, p.prio - Math.floor(waited/aging));
        }

        while(st.some(p=>p.remaining>0) && guard<200000){
            guard++;
            while(i<st.length && st[i].arrival<=t){ 
                if(!queued.has(st[i].name) && st[i].remaining>0 && st[i]!==run){ 
                    q.push(st[i]); 
                    queued.add(st[i].name); 
                    waitSince[st[i].name]=t;
                } 
                i++; 
            }

            if(!run){
                if(q.length===0){ 
                    t++; 
                    continue; 
                }
                q.sort((a,b)=> effPrio(a)-effPrio(b) || a.arrival-b.arrival);
                run=q.shift(); queued.delete(run.name); 
                waitSince[run.name]=null; 
                runLeft=Math.min(quantum, run.remaining);
            }
            if(run.start===null) 
                run.start=t;
            pushTick(run.name);
            run.remaining--; 
            runLeft--; 
            t++;
            while(i<st.length && st[i].arrival<=t){ 
                if(!queued.has(st[i].name) && st[i].remaining>0 && st[i]!==run){ 
                    q.push(st[i]); 
                    queued.add(st[i].name); 
                    waitSince[st[i].name]=t;
                } 
                i++; 
            }
            if(run.remaining<=0){ 
                run.finish=t; 
                run=null; 
            }
            else if(runLeft<=0){ 
                q.push(run); 
                queued.add(run.name); 
                waitSince[run.name]=t; 
                run=null; 
            }
        }
        closeSeg();
    }

    return {intervals, stats:st, contextSwitches};
}

let SIM=null, playTimer=null, curT=0;

function run(){
  const procs=readProcs();
  const algo=document.getElementById('algo').value;
  const quantum=+document.getElementById('quantum').value||2;
  const aging=+document.getElementById('aging').value||3;
  const {intervals, stats, contextSwitches}=simulate(procs, algo, quantum, aging);
  const maxT=Math.max(...intervals.map(iv=>iv.end), ...stats.map(s=>s.finish||0), 1);
  const scale=Math.max(28, Math.min(60, 640/maxT));
  SIM={procs, intervals, stats, maxT, scale, contextSwitches};
  curT=maxT;
  stopPlay();
  renderMetrics();
  renderAtTime(curT);
  document.getElementById('resultCard').style.display='block';
}

function renderAtTime(t){
    const {procs, intervals, stats, maxT, scale}=SIM;
    const inner=document.getElementById('ganttInner');
    inner.innerHTML='';
    procs.forEach(p=>{
        const row=document.createElement('div'); row.className='row';
        const lbl=document.createElement('div'); lbl.className='rowlabel'; lbl.textContent=p.name;
        const track=document.createElement('div'); track.className='track'; track.style.width=(maxT*scale)+'px';
        intervals.filter(iv=>iv.name===p.name).forEach(iv=>{
            const end=Math.min(iv.end, t);

            if(end<=iv.start) 
                return;

            const b=document.createElement('div'); 
            b.className='bar';
            b.style.left=(iv.start*scale)+'px'; 
            b.style.width=((end-iv.start)*scale-2)+'px';
            b.style.background=`var(${p.color[0]})`; 
            b.style.borderColor=`var(${p.color[1]})`;
            track.appendChild(b);
        });

        const s=stats.find(x=>x.name===p.name);
        if(s.arrival<=t){
            const dA=document.createElement('div'); 
            dA.className='diamond'; 
            dA.style.left=(s.arrival*scale)+'px';
            dA.style.borderColor=`var(${p.color[1]})`; 
            dA.style.color=`var(${p.color[1]})`;
            track.appendChild(dA);
        }

        if(s.finish!==null && s.finish<=t){
        const dF=document.createElement('div'); 
        dF.className='diamond done'; 
        dF.style.left=(s.finish*scale)+'px';
        dF.style.borderColor=`var(${p.color[1]})`; 
        dF.style.color=`var(${p.color[1]})`;
        track.appendChild(dF);
        }

        if(t<maxT){
        const cur=document.createElement('div'); 
        cur.className='cursor'; 
        cur.style.left=(t*scale)+'px';
        track.appendChild(cur);
        }

        row.append(lbl,track); 
        inner.appendChild(row);
    });

    const axis=document.createElement('div'); 
    axis.className='axis';
    const albl=document.createElement('div'); 
    albl.className='rowlabel';
    const atrack=document.createElement('div'); 
    atrack.className='track'; 
    atrack.style.width=(maxT*scale)+'px';
    for(let x=0;x<=maxT;x++){ 
        const tk=document.createElement('div'); 
        tk.className='tick'; 
        tk.style.left=(x*scale)+'px'; 
        tk.textContent=x; 
        atrack.appendChild(tk); 
    }
    axis.append(albl,atrack); 
    inner.appendChild(axis);

    //legenda do grafico
    const legend=document.getElementById('legend'); 
    legend.innerHTML='';
    legend.innerHTML=`
        <span>
            <span class="diamond" style="position:static;transform:rotate(45deg);width:8px;height:8px;border-color:var(--muted);color:var(--muted)">
            </span>Chegada</span>
            <span><span class="diamond done" style="position:static;transform:rotate(45deg);width:8px;height:8px;border-color:var(--muted);color:var(--muted);background:var(--muted)">
            </span>Conclusão</span>`;

    document.getElementById('clockLine').textContent = t>=maxT ? `t = ${maxT} (Concluído)` : `t = ${t}`;
}

function renderMetrics(){
    const {stats, contextSwitches}=SIM;
    const tb=document.getElementById('metricsTable');
    let waitSum=0, turnSum=0, respSum=0;
    let html='<tr><th>Processo</th><th>Chegada</th><th>Duração</th><th>Início</th><th>Término</th><th>Espera</th><th>Retorno</th><th>Resposta</th></tr>';
    stats.forEach(s=>{
        const turn=s.finish-s.arrival, wait=turn-s.burst, resp=s.start-s.arrival;
        waitSum+=wait; 
        turnSum+=turn;
        respSum+=resp;
        html+=`<tr><td>${s.name}</td><td>${s.arrival}</td><td>${s.burst}</td><td>${s.start}</td><td>${s.finish}</td><td>${wait}</td><td>${turn}</td><td>${resp}</td></tr>`;
    });
    const n=stats.length;
    html+=`<tr><td colspan="5">Médias</td><td>${(waitSum/n).toFixed(2)}</td><td>${(turnSum/n).toFixed(2)}</td><td>${(respSum/n).toFixed(2)}</td></tr>`;
    html+=`<tr><td colspan="8" style="text-align:center;font-weight:600;">Trocas de contexto: ${contextSwitches}</td></tr>`;
    tb.innerHTML=html;
}

function togglePlay(){
    if(!SIM) 
        return;
    if(playTimer){ stopPlay(); 
        return; 
    }
    if(curT>=SIM.maxT) 
        curT=0;

    const btn=document.getElementById('playBtn'); 
    btn.textContent='⏸ Pausar';
    const speed=+document.getElementById('speed').value;
    playTimer=setInterval(()=>{
        curT++;
        renderAtTime(curT);
        if(curT>=SIM.maxT) stopPlay();
    }, speed);
}


function stopPlay(){
    if(playTimer){ clearInterval(playTimer); 
        playTimer=null; 
        }

    const btn=document.getElementById('playBtn'); 
    if(btn) 
        btn.textContent='▶ Play';
}


function resetPlay(){
    if(!SIM) 
        return;

    stopPlay(); 
    curT=0; 
    renderAtTime(curT);
}


document.getElementById('speed').addEventListener('change', ()=>{
    if(playTimer){ 
        stopPlay(); 
        togglePlay(); 
    }
});


run();