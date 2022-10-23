class SlideRuler {
    constructor(elementID,config={
      width:'300px',
      height:'70px',
      from: -10,
      to: 10,
      step: 5,
      backColor: "black",
      foreColor: "white",
      pointer:"white",
      onchange:(v)=>{ console.log(`slideRuler onchanged : ${v}`) }
    }){
      this.element = document.getElementById(elementID)
      
      this.backColor = config.backColor?config.backColor:"black"
      this.foreColor = config.foreColor?config.foreColor:"white"
      this.pointer = config.pointer?config.pointer:this.foreColor
      this.cursor = config.cursor?config.cursor:"default"
      
      this.element.style.backgroundColor=this.backColor
      this.element.style.overflow="hidden"
      this.element.style.userSelect='none'
      this.element.style.cursor= this.cursor
      
      var sign = config.sign?config.sign:""
      
      this.name = "SlideRuler" + (new Date().getTime())
      
      this.components={
        tick:function(value,isKeyTick,tickStyle){
          const style=`
            <div
              style="
                width:${tickStyle.width||10}px;
                height:100%;
                display:flex;
                float:left;
                flex: 1;
                flex-direction:column;
                color:${tickStyle.tickColor};
                align-items: center;
              "
            >
              <div style="width:${isKeyTick?5:3}px;height:${isKeyTick?14:8}px;">
                <svg style="width:100%;">
                  <line x1="0" y1="0" x2="0" y2="${isKeyTick?14:8}" style="stroke:${tickStyle.tickColor};stroke-width:${isKeyTick?5:3}px"/> 
                </svg>
              </div>
              <div style="width:100%;height:auto;flex:1;display:flex;justify-content:center;align-items:center">
                <span style="${isKeyTick?'':'display:none;'}">
                  ${value||0}${sign}
                </span>
              </div>
              <div style="width:${isKeyTick?5:3}px;height:12px;justify-content:center;align-items:end">
                <svg style="width:100%;height:${isKeyTick?14:8}px;${tickStyle.bothSide||'display:none;'}">
                  <line x1="0" y1="0" x2="0" y2="${isKeyTick?14:8}" style="stroke: ${tickStyle.tickColor};stroke-width:${isKeyTick?5:3}px"/> 
                </svg>
              </div>
            </div>
          `
          
          return style
        },
        pointer:function(isKeyTick){
          
        },
      }
      
      this.element.style.width = config.width||'300px'
      this.element.style.height = config.height||'70px'
      this.ticks = ""
      
      config.from = config.from?config.from:-10
      config.to = config.to?config.to:10
      config.step = config.step?config.step:5
      
      if(config.tickStyle){
        config.tickStyle.tickColor = config.tickStyle.tickColor?config.tickStyle.tickColor:this.foreColor
      }else{
        config.tickStyle = { tickColor:this.foreColor }
      }
      
      for(let i=config.from;i<=config.to;i++){
        let isKeyTick = (i%config.step==0)
        this.ticks+= this.components.tick(i,isKeyTick,config.tickStyle)
      }
      
      const elWidth = 10*(Math.abs(config.from-config.to)+1)
      
      this.element.innerHTML = `
        <div style="width:100%;height:100%;position:relative;">
          <div
            id="${this.name}-ticker"
            style="
              width:${elWidth}px;
              height:100%;
              background-color:${this.backColor};
              position:relative;
              left:50%;
              margin-left:-${5*(Math.abs(config.from-config.to)+1)}px;
            "
          >
            ${this.ticks}
          </div>
          <div style="position:absolute;width:100%;height:100%;top:0;text-align:center">
            <div style="height:100%;padding:0 5px 0 5px;margin:0px auto;background:linear-gradient(to right,transparent, ${this.backColor} , ${this.backColor} , ${this.backColor} ,transparent);display:inline-block;">
              <div style="height:100%;display:flex;flex-direction:column;align-items:center;">
                <div style="width:20px;height:12px;">
                  <svg style="width:20px;height:12px;">
                    <polygon fill="${this.pointer}" points="2,0 18,0 10,12"/>
                  </svg>
                </div>
                
                <div style="flex:1 1 auto;font-size:24px;color:${this.foreColor};display:flex;justify-content:center;align-items:center;">
                  <span id="${this.name}-tickValue" >0</span>${sign}
                </div>
                
                <div style="width:20px;height:12px;">
                  <svg style="width:20px;height:12px;display:none;">
                    <polygon fill="${this.pointer}" points="2,0 18,0 10,12"/>
                  </svg>
                </div>
              </div>
            </div>
          </div>
        </div>
      `
      
      var valueWindow = document.getElementById(this.name+"-tickValue")
      var ticker = document.getElementById(this.name+"-ticker")
      var isDrawing = false
      var x = 0
      var deltaX = 0
      
      var onchange = config.onchange?config.onchange:null
      
      
      const writeX = _x => {
        x = _x
        isDrawing = true
      }
      const moving = xValue => {
        console.log(xValue + " ==== " + x)
        deltaX += (parseInt(xValue)-parseInt(x))
        x = xValue
  
        if(deltaX<-elWidth/2){
          deltaX=-elWidth/2
        }else if(deltaX>elWidth/2){
          deltaX=elWidth/2
        }
        ticker.style.left=`calc(50% + ${deltaX}px)`
  
        const value = Math.floor(-deltaX/10)
        valueWindow.innerText = value;
  
        if(onchange){
          onchange(value)
        }
      }
      
      const pDown = e => {
        //e.preventDefault()
        console.log(e.type)
        writeX(e.screenX)
      }
      const tStart = e => {
        //e.preventDefault()
        const touches = e.changedTouches
        console.log(e.type)
        writeX(touches[0].pageX)
      }
      const pMove = e => {
        if (isDrawing === true) {
          console.log(e.screenX + " : " + x)
          moving(e.screenX)
        }
      }
      const tMove = e => {
        if (isDrawing === true) {
          const touches = e.changedTouches
          console.log(touches[0].pageX + " : " + x)
          moving(touches[0].pageX)
        }
      }
      const srUp = e => {
        isDrawing = false
      }
      
      this.element.addEventListener('pointerdown', pDown)
      //this.element.addEventListener('touchstart', tStart)
      
      window.addEventListener('pointermove', pMove)
      //window.addEventListener('touchmove', tMove)
      
      window.addEventListener('pointerup', srUp)
      //window.addEventListener('touchend', srUp)
    }
  }
  
  const slideRuler = new SlideRuler("slideRuler1",{
    width: '350px',
    height: '60px',
    from:-45,
    to:45,
    step:5,
    sign:"°",
    foreColor:"lightgray",
    pointer:"#dd0000",
    cursor:"w-resize",
    onchange:(v)=>{
      const image = document.getElementById("image")
      image.style.transform = `rotate(${v}deg)`
    }
  })