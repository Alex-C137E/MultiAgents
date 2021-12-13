package main

import(
  agents "gitlab.utc.fr/jedescam/projet_ia04/agents";
  //milieu "gitlab.utc.fr/jedescam/projet_ia04/milieu";
  gg "github.com/fogleman/gg";
  "math/rand";
  "math";
  "fmt";
  "time";
)

func remove(slice []*agents.Boid	, s int) []*agents.Boid {
    return append(slice[:s], slice[s+1:]...)
}


func main(){
  // initialiser à la taille de la fenetre
  var monde agents.Monde
  monde.Agents = make([][][]*agents.Boid,1000)
  monde.Preds = make([][][]*agents.Pred,1000)
  for i := 0; i < len(monde.Agents); i++ {
		monde.Agents[i] = make([][]*agents.Boid, 1000)
    monde.Preds[i] = make([][]*agents.Pred, 1000)
		/*for j := 0; j < len(monde.Agents[i]); j++ {
      monde.Agents[i][j] = make([]*agents.Boid, 20)
      monde.Preds[i][j] = make([]*agents.Pred, 20)
		}*/
	}

  // init sources
  for s:=0; s<2; s++{
    sx := 500 +150*(-1.0 + rand.Float64()*2)
    sy := 500 +200*(-1.0 + rand.Float64()*2)
    source := agents.New_source(0,[2]float64{sx,sy})
    monde.Sources = append(monde.Sources, &source)
  }
  // init boids
  for i:=0 ; i<50; i++{
			bx := 500 + 70*(-1.0 + rand.Float64()*2)
			by := 500 + 70*(-1.0 + rand.Float64()*2)
			boid := agents.New_boid(i,[2]float64{bx,by}, 1) // 20 c'est le rayon du flock ici !
			monde.Agents[int(math.Trunc(bx))][int(math.Trunc(by))] = append(monde.Agents[int(math.Trunc(bx))][int(math.Trunc(by))], &boid)
		}
  // init pred
  for i:=0 ; i<40; i++{
			px := 250 + 70*(-1.0 + rand.Float64()*2)
			py := 200 + 70*(-1.0 + rand.Float64()*2)
			pred := agents.New_pred([2]float64{px,py}, 1) // 20 c'est le rayon du flock ici !
			monde.Preds[int(math.Trunc(px))][int(math.Trunc(py))] = append(monde.Preds[int(math.Trunc(px))][int(math.Trunc(py))], &pred)
		}
  // --------- Fin init ------------

  // --------- Actualisation du Monde ------------
  ticker := time.NewTicker(20 * time.Millisecond)
	done := make(chan bool)
	c := 0
	go func() {
	   for {
	      select {
	       case <-done:
	          return
	       case t := <-ticker.C:
               //compteur := 0
				       fmt.Println("Tick at", t)
							 dc := gg.NewContext(1000, 1000)
               monde.Sources[0].Draw(dc)
               monde.Sources[1].Draw(dc)
               compte := 0
               for i,_ := range(monde.Agents){
                 for j,_ := range(monde.Agents[i]){
                   if len(monde.Agents[i][j]) == 0 {
                     continue
                   }
                   boids := monde.Agents[i][j]
                   if len(monde.Preds[i][j]) != 0{ // on actualise pas ça [i][j]
                     fmt.Println("aka")
                     monde.Agents[i][j] = []*agents.Boid{}
                     continue
                   }
                   compte += len(boids)
                   c:=0
                   for _, boid := range(boids){
                     boid.Get_couches(&monde)
                     boid.Flock()

                     c+=1
                   }
               }
             }
             fmt.Println("nombre de boids", compte)

              for i,_ := range(monde.Agents){
                 for j,_ := range(monde.Agents[i]){
                   if len(monde.Agents[i][j]) == 0 {
                     continue
                   }
                   boids := monde.Agents[i][j]
                   for _, boid := range(boids){
                     boid.Update() // ajouter un retour ici ! -> choper les positions pour actualiser le tableau
                     boid.Border(1000.0,1000.0)
                     boid.Draw(dc)

                   }

                 }
               }
               for i,_ := range(monde.Preds){
                  for j,_ := range(monde.Preds[i]){
                    if len(monde.Preds[i][j]) == 0 {
                      continue
                    }
                    boids := monde.Preds[i][j]
                    for _, boid := range(boids){
                      boid.Update()
                      boid.Draw(dc)
                    }
                  }
                }
							 /*for _,boid := range(monde.Agents){

							 }
							 for _, boid := range(monde.Agents){
                  boid.Update()
									if !(boid.Is_alive()){
										monde.Agents = remove(monde.Agents, compteur)
									}*/
								img := fmt.Sprintf("out_n%v.png",c)
								dc.Fill()
								dc.SavePNG(img)
								c+=1
								// travailler localement avec les boid -> créer des sous-monde.boidss.
	        }
				}
	    }()
	  time.Sleep(80000 * time.Millisecond)
	  ticker.Stop()
		done <- true
		fmt.Println("Ticker stopped")
}
