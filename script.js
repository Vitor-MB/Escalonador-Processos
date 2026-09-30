const COLORS=[['--p1','--p1d'],['--p2','--p2d'],['--p3','--p3d'],['--p4','--p4d'],['--p5','--p5d'],['--p6','--p6d']];

// Endpoint base da API (não altere se o backend estiver em outra porta)
const API_BASE = 'http://localhost:8080';
// Mapa de metadados dos algoritmos retornados pela API
const metadataMap = {};

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
    const meta = metadataMap[a];
    const usesQuantum = meta ? meta.uses_quantum : (a==='rr'||a==='rr-prio-aging');
    const usesAging = meta ? meta.uses_aging : (a==='rr-prio-aging');
    document.getElementById('quantumBox').classList.toggle('show', usesQuantum);
    document.getElementById('agingBox').classList.toggle('show', usesAging);
}

// Inicializa a UI buscando os algoritmos no backend (se disponível)
async function init(){
    try{
        const resp = await fetch(API_BASE + '/algorithm');
        if(resp.ok){
            const data = await resp.json();
            const sel = document.getElementById('algo');
            sel.innerHTML = data.algorithms.map(a=>`<option value="${a.algorithm}">${a.algorithm}</option>`).join('');
            data.algorithms.forEach(a=>metadataMap[a.algorithm]=a);
            syncAlgoUI();
            hideError();
        }
    }catch(e){
        console.warn('Não foi possível buscar algoritmos do backend:', e);
        showError('Não foi possível carregar algoritmos do backend. Verifique se o servidor está rodando.');
    }
}
document.addEventListener('DOMContentLoaded', init);

// Adiciona alguns Elementos de exemplo para os processos
[{a:0,b:5,p:2},{a:1,b:3,p:1},{a:2,b:8,p:3},{a:3,b:6,p:2},{a:4,b:2,p:1}].forEach(addRow);
syncAlgoUI();

function readProcs(){
  return [...document.querySelectorAll('#procBody tr')].map((tr,i)=>{
    const g=f=>tr.querySelector(`[data-f="${f}"]`).value;
    return {name:g('name')||('P'+(i+1)), a:+g('a'), b:+g('b'), p:+g('p'), color:COLORS[i%COLORS.length]};
  });
}

// REMOVIDO: simulador local removido para depender exclusivamente do backend.

// Mostrar/ocultar banner de erro
function showError(msg){
    const b=document.getElementById('errorBanner');
    if(b){ b.textContent=msg; b.style.display='block'; }
}

function hideError(){
    const b=document.getElementById('errorBanner');
    if(b){ b.textContent=''; b.style.display='none'; }
}

let SIM=null, playTimer=null, curT=0;

function run(){
    const procs=readProcs();
    const algo=document.getElementById('algo').value;
    const quantum=+document.getElementById('quantum').value||2;
    const aging=+document.getElementById('aging').value||3;

    // Chama o backend; se falhar, exibe erro para o usuário
    (async ()=>{
        try{
            hideError();
            const body = {
                algorithm: algo,
                quantum: quantum,
                aging: aging,
                processes: procs.map(p=>({name:p.name, arrival:p.a, burst:p.b, priority:p.p}))
            };
            const resp = await fetch(API_BASE + '/simulate', {
                method: 'POST',
                headers: {'Content-Type':'application/json'},
                body: JSON.stringify(body)
            });

            if(!resp.ok){
                showError(`Falha na requisição: ${resp.status} ${resp.statusText}`);
                return;
            }

            const data = await resp.json();
            // Converte a resposta do backend para a estrutura usada pelo front
            const intervals = data.intervals.map(i=>({name:i.process_name, start:i.start, end:i.finish}));
            const stats = data.processes.map(p=>({name:p.name, arrival:p.arrival, burst:p.burst, start:p.start, finish:p.finish}));
            const maxT = data.total_time || Math.max(...intervals.map(iv=>iv.end), ...stats.map(s=>s.finish||0),1);
            const scale = Math.max(28, Math.min(60, 640/maxT));
            SIM={procs, intervals, stats, maxT, scale, contextSwitches: data.context_switches||0};
            curT=maxT;
            stopPlay();
            renderMetrics();
            renderAtTime(curT);
            document.getElementById('resultCard').style.display='block';
        }catch(e){
            console.error('Erro ao chamar backend:', e);
            showError('Erro ao comunicar com o backend. Verifique se o servidor está rodando em ' + API_BASE);
            return;
        }
    })();
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